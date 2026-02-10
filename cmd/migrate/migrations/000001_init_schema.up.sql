

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


CREATE TABLE IF NOT EXISTS accounts (
  id int(11) NOT NULL AUTO_INCREMENT,
  account_number int(11) NOT NULL,
  account_name varchar(50) NOT NULL,
  account_type int(11) NOT NULL,
  building_id int(11) NOT NULL,
  isDefault tinyint(1) NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp(),
  updated_at datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY acc_acc_type_fk (account_type),
  KEY accounts_building_fk (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table account_types
--

CREATE TABLE IF NOT EXISTS account_types (
  id int(11) NOT NULL AUTO_INCREMENT,
  typeName varchar(250) NOT NULL,
  type varchar(20) NOT NULL,
  sub_type varchar(20) NOT NULL,
  typeStatus varchar(10) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table bills
--

CREATE TABLE IF NOT EXISTS bills (
  id int(11) NOT NULL AUTO_INCREMENT,
  bill_no varchar(255) NOT NULL,
  transaction_id int(11) NOT NULL,
  bill_date date NOT NULL,
  due_date date NOT NULL,
  ap_account_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  user_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text NOT NULL,
  cancel_reason text DEFAULT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  building_id int(11) NOT NULL,
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY bill_unit_id (unit_id),
  KEY bill_people_id (people_id),
  KEY bill_user_id (user_id),
  KEY bill_building_id (building_id),
  KEY fk_bill_account_id (ap_account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table bill_expense_lines
--

CREATE TABLE IF NOT EXISTS bill_expense_lines (
  id int(11) NOT NULL AUTO_INCREMENT,
  bill_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  description text DEFAULT NULL,
  amount decimal(10,2) NOT NULL,
  PRIMARY KEY (id),
  KEY bill_id (bill_id),
  KEY account_id (account_id),
  KEY unit_id (unit_id),
  KEY people_id (people_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table bill_payments
--

CREATE TABLE IF NOT EXISTS bill_payments (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  reference varchar(255) NOT NULL,
  date date NOT NULL,
  bill_id int(11) NOT NULL,
  user_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY bp_transaction_id (transaction_id),
  KEY bp_bill_id (bill_id),
  KEY bp_user_id (user_id),
  KEY bp_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table buildings
--

CREATE TABLE IF NOT EXISTS buildings (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(250) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  UNIQUE KEY name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table checks
--

CREATE TABLE IF NOT EXISTS checks (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  check_date date NOT NULL,
  reference_number varchar(50) DEFAULT NULL,
  payment_account_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  memo text DEFAULT NULL,
  total_amount decimal(10,2) DEFAULT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY transaction_id (transaction_id),
  KEY payment_account_id (payment_account_id),
  KEY building_id (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table credit_memo
--

CREATE TABLE IF NOT EXISTS credit_memo (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  reference varchar(255) NOT NULL,
  date date NOT NULL,
  user_id int(11) NOT NULL,
  deposit_to int(11) NOT NULL,
  liability_account int(11) NOT NULL,
  people_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  unit_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY cm_transaction_id (transaction_id),
  KEY cm_user_id (user_id),
  KEY cm_deposit_to (deposit_to),
  KEY cm_liability_account (liability_account),
  KEY cm_people_id (people_id),
  KEY cm_building_id (building_id),
  KEY cm_unit_id (unit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table expense_lines
--

CREATE TABLE IF NOT EXISTS expense_lines (
  id int(11) NOT NULL AUTO_INCREMENT,
  check_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  description text DEFAULT NULL,
  amount decimal(10,2) NOT NULL,
  PRIMARY KEY (id),
  KEY check_id (check_id),
  KEY account_id (account_id),
  KEY unit_id (unit_id),
  KEY people_id (people_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table invoices
--

CREATE TABLE IF NOT EXISTS invoices (
  id int(11) NOT NULL AUTO_INCREMENT,
  invoice_no varchar(255) NOT NULL,
  transaction_id int(11) NOT NULL,
  sales_date date NOT NULL,
  due_date date NOT NULL,
  ar_account_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  user_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text NOT NULL,
  cancel_reason text DEFAULT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  building_id int(11) NOT NULL,
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  amount_cents bigint(20) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY invoice_unit_id (unit_id),
  KEY invoice_people_id (people_id),
  KEY invoice_user_id (user_id),
  KEY invoice_building_id (building_id),
  KEY fk_invoice_account_id (ar_account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table invoice_applied_credits
--

CREATE TABLE IF NOT EXISTS invoice_applied_credits (
  id int(11) NOT NULL AUTO_INCREMENT,
  invoice_id int(11) NOT NULL,
  credit_memo_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text NOT NULL,
  date date NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY iac_invoice_id (invoice_id),
  KEY iac_credit_memo_id (credit_memo_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table invoice_applied_discounts
--

CREATE TABLE IF NOT EXISTS invoice_applied_discounts (
  id int(11) NOT NULL AUTO_INCREMENT,
  reference varchar(255) NOT NULL,
  invoice_id int(11) NOT NULL,
  transaction_id int(11) NOT NULL,
  ar_account int(11) NOT NULL,
  income_account int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text NOT NULL,
  date date NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY iad_invoice_id (invoice_id),
  KEY iad_transaction_id (transaction_id),
  KEY iad_ar_account (ar_account),
  KEY iad_income_account (income_account)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table invoice_items
--

CREATE TABLE IF NOT EXISTS invoice_items (
  id int(11) NOT NULL AUTO_INCREMENT,
  invoice_id int(11) NOT NULL,
  item_id int(11) NOT NULL,
  item_name varchar(250) NOT NULL,
  previous_value decimal(10,3) DEFAULT NULL,
  current_value decimal(10,3) DEFAULT NULL,
  qty decimal(10,3) DEFAULT NULL,
  rate varchar(100) DEFAULT NULL,
  total decimal(10,2) NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  qty_scaled bigint(20) DEFAULT NULL,
  rate_scaled bigint(20) DEFAULT NULL,
  total_cents bigint(20) DEFAULT NULL,
  previous_value_cents bigint(20) DEFAULT NULL,
  current_value_cents bigint(20) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY invoice_items_inv_id (invoice_id),
  KEY invoice_items_item_id (item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table invoice_payments
--

CREATE TABLE IF NOT EXISTS invoice_payments (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  reference varchar(255) NOT NULL,
  date date NOT NULL,
  invoice_id int(11) NOT NULL,
  user_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY ip_transaction_id (transaction_id),
  KEY ip_invoice_id (invoice_id),
  KEY ip_user_id (user_id),
  KEY ip_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table items
--

CREATE TABLE IF NOT EXISTS items (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(250) NOT NULL,
  type enum('inventory','non inventory','service','discount','payment') NOT NULL,
  description text NOT NULL,
  asset_account int(11) DEFAULT NULL,
  income_account int(11) DEFAULT NULL,
  cogs_account int(11) DEFAULT NULL,
  expense_account int(11) DEFAULT NULL,
  on_hand decimal(10,2) NOT NULL,
  avg_cost decimal(10,2) NOT NULL,
  date date NOT NULL,
  building_id int(11) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY item_asset_account_pk (asset_account),
  KEY item_income_account_pk (income_account),
  KEY item_cogs_account_pk (cogs_account),
  KEY item_expense_account_pk (expense_account),
  KEY item_building_fk (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table journal
--

CREATE TABLE IF NOT EXISTS journal (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  reference varchar(255) NOT NULL,
  journal_date date NOT NULL,
  building_id int(11) NOT NULL,
  memo text DEFAULT NULL,
  total_amount decimal(10,2) DEFAULT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY transaction_id (transaction_id),
  KEY building_id (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table journal_lines
--

CREATE TABLE IF NOT EXISTS journal_lines (
  id int(11) NOT NULL AUTO_INCREMENT,
  journal_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  description text DEFAULT NULL,
  debit decimal(10,2) DEFAULT 0.00,
  credit decimal(10,2) DEFAULT 0.00,
  PRIMARY KEY (id),
  KEY journal_id (journal_id),
  KEY account_id (account_id),
  KEY unit_id (unit_id),
  KEY people_id (people_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table leases
--

CREATE TABLE IF NOT EXISTS leases (
  id int(11) NOT NULL AUTO_INCREMENT,
  people_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  unit_id int(11) NOT NULL,
  start_date date NOT NULL,
  end_date date DEFAULT NULL,
  rent_amount decimal(10,2) NOT NULL,
  deposit_amount decimal(10,2) NOT NULL,
  service_amount decimal(10,2) NOT NULL,
  lease_terms text NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  PRIMARY KEY (id),
  KEY idx_people_id (people_id),
  KEY idx_building_id (building_id),
  KEY idx_unit_id (unit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table lease_files
--

CREATE TABLE IF NOT EXISTS lease_files (
  id int(11) NOT NULL AUTO_INCREMENT,
  lease_id int(11) NOT NULL,
  filename varchar(255) NOT NULL,
  original_name varchar(255) NOT NULL,
  file_path varchar(500) NOT NULL,
  file_type varchar(100) NOT NULL,
  file_size int(11) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY idx_lease_id (lease_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table people
--

CREATE TABLE IF NOT EXISTS people (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  phone varchar(20) NOT NULL,
  type_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY people_type_id_fk (type_id),
  KEY people_building_id_fk (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table people_types
--

CREATE TABLE IF NOT EXISTS people_types (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(50) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table periods
--

CREATE TABLE IF NOT EXISTS periods (
  id int(11) NOT NULL AUTO_INCREMENT,
  period_name varchar(50) NOT NULL,
  start date NOT NULL,
  end date NOT NULL,
  building_id int(11) NOT NULL,
  is_closed int(2) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY period_building_id (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table permissions
--

CREATE TABLE IF NOT EXISTS permissions (
  id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  module varchar(50) NOT NULL,
  action varchar(50) NOT NULL,
  `key` varchar(100) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table readings
--

CREATE TABLE IF NOT EXISTS readings (
  id int(11) NOT NULL AUTO_INCREMENT,
  item_id int(11) NOT NULL,
  unit_id int(11) NOT NULL,
  lease_id int(11) DEFAULT NULL,
  reading_month varchar(10) DEFAULT NULL,
  reading_year varchar(5) DEFAULT NULL,
  reading_date date NOT NULL,
  previous_value decimal(10,3) DEFAULT NULL,
  current_value decimal(10,3) DEFAULT NULL,
  unit_price decimal(10,2) DEFAULT NULL,
  total_amount decimal(10,2) DEFAULT NULL,
  notes text DEFAULT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY idx_item_id (item_id),
  KEY idx_unit_id (unit_id),
  KEY idx_lease_id (lease_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table receipt_items
--

CREATE TABLE IF NOT EXISTS receipt_items (
  id int(11) NOT NULL AUTO_INCREMENT,
  receipt_id int(11) NOT NULL,
  item_id int(11) NOT NULL,
  item_name varchar(250) NOT NULL,
  previous_value decimal(10,3) DEFAULT NULL,
  current_value decimal(10,3) DEFAULT NULL,
  qty decimal(10,2) DEFAULT NULL,
  rate varchar(100) DEFAULT NULL,
  total decimal(10,2) NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY sri_receipt_id (receipt_id),
  KEY sri_item_id (item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table roles
--

CREATE TABLE IF NOT EXISTS roles (
  id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  owner_user_id int(11) NOT NULL,
  name varchar(100) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  UNIQUE KEY uq_owner_role_name (owner_user_id,name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table role_permissions
--

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id bigint(20) UNSIGNED NOT NULL,
  permission_id bigint(20) UNSIGNED NOT NULL,
  PRIMARY KEY (role_id,permission_id),
  KEY fk_rp_permission (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table sales_receipt
--

CREATE TABLE IF NOT EXISTS sales_receipt (
  id int(11) NOT NULL AUTO_INCREMENT,
  receipt_no int(11) NOT NULL,
  transaction_id int(11) NOT NULL,
  receipt_date date NOT NULL,
  unit_id int(11) DEFAULT NULL,
  people_id int(11) DEFAULT NULL,
  user_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  amount decimal(10,2) NOT NULL,
  description text DEFAULT NULL,
  cancel_reason text DEFAULT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  building_id int(11) NOT NULL,
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id),
  KEY sr_unit_id (unit_id),
  KEY sr_people_id (people_id),
  KEY sr_user_id (user_id),
  KEY sr_building_id (building_id),
  KEY sr_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Table structure for table splits
--

CREATE TABLE IF NOT EXISTS splits (
  id int(11) NOT NULL AUTO_INCREMENT,
  transaction_id int(11) NOT NULL,
  account_id int(11) NOT NULL,
  people_id int(11) DEFAULT NULL,
  unit_id int(11) DEFAULT NULL,
  debit decimal(10,2) DEFAULT NULL,
  credit decimal(10,2) DEFAULT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at datetime NOT NULL DEFAULT current_timestamp(),
  updated_at datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  debit_cents bigint(20) DEFAULT NULL,
  credit_cents bigint(20) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY td_transaction_id (transaction_id),
  KEY td_account_id (account_id),
  KEY td_people_id (people_id),
  KEY fk_splits_unit (unit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table transactions
--

CREATE TABLE IF NOT EXISTS transactions (
  id int(11) NOT NULL AUTO_INCREMENT,
  type enum('invoice','payment','check','deposit','bill','credit memo','sales receipt','journal','bill credit','bill payment','credit applied') NOT NULL,
  transaction_date date NOT NULL,
  transaction_number varchar(255) NOT NULL,
  memo text NOT NULL,
  status enum('0','1') NOT NULL DEFAULT '1',
  created_at datetime NOT NULL DEFAULT current_timestamp(),
  updated_at datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  building_id int(11) NOT NULL,
  user_id int(11) NOT NULL,
  unit_id int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY transaction_type_fk (type),
  KEY transaction_building_fk (building_id),
  KEY transaction_user_fk (user_id),
  KEY transaction_unit_fk (unit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table units
--

CREATE TABLE IF NOT EXISTS units (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(20) NOT NULL,
  building_id int(11) NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp(),
  updated_at timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY unit_building_id_fk (building_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table users
--

CREATE TABLE IF NOT EXISTS users (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  username varchar(20) NOT NULL,
  phone varchar(20) NOT NULL,
  password varchar(10) NOT NULL,
  parent_user_id int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY parent_user_id (parent_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table users_building
--

CREATE TABLE IF NOT EXISTS users_building (
  user_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  KEY user_building_building_id (building_id),
  KEY user_building_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table user_building_roles
--

CREATE TABLE IF NOT EXISTS user_building_roles (
  user_id int(11) NOT NULL,
  building_id int(11) NOT NULL,
  role_id bigint(20) UNSIGNED NOT NULL,
  PRIMARY KEY (user_id,building_id,role_id) USING BTREE,
  KEY fk_ubr_building (building_id),
  KEY fk_ubr_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Constraints for dumped tables
--

--
-- Constraints for table accounts
--
ALTER TABLE accounts
  ADD CONSTRAINT acc_acc_type_fk FOREIGN KEY (account_type) REFERENCES account_types (id),
  ADD CONSTRAINT accounts_building_fk FOREIGN KEY (building_id) REFERENCES buildings (id);

--
-- Constraints for table bill_payments
--
ALTER TABLE bill_payments
  ADD CONSTRAINT fk_bill_payments_account FOREIGN KEY (account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_bill_payments_bill FOREIGN KEY (bill_id) REFERENCES bills (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_bill_payments_transaction FOREIGN KEY (transaction_id) REFERENCES `transactions` (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table checks
--
ALTER TABLE checks
  ADD CONSTRAINT checks_ibfk_1 FOREIGN KEY (transaction_id) REFERENCES `transactions` (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT checks_ibfk_2 FOREIGN KEY (payment_account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT checks_ibfk_3 FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table credit_memo
--
ALTER TABLE credit_memo
  ADD CONSTRAINT fk_cm_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_deposit_to FOREIGN KEY (deposit_to) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_liability_account FOREIGN KEY (liability_account) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_people FOREIGN KEY (people_id) REFERENCES people (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_transaction FOREIGN KEY (transaction_id) REFERENCES `transactions` (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_cm_user FOREIGN KEY (user_id) REFERENCES `users` (id) ON UPDATE CASCADE;

--
-- Constraints for table expense_lines
--
ALTER TABLE expense_lines
  ADD CONSTRAINT expense_lines_ibfk_1 FOREIGN KEY (check_id) REFERENCES `checks` (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT expense_lines_ibfk_2 FOREIGN KEY (account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT expense_lines_ibfk_3 FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT expense_lines_ibfk_4 FOREIGN KEY (people_id) REFERENCES people (id) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table invoices
--
ALTER TABLE invoices
  ADD CONSTRAINT fk_invoice_account_id FOREIGN KEY (ar_account_id) REFERENCES `accounts` (id),
  ADD CONSTRAINT fk_invoice_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_invoice_people FOREIGN KEY (people_id) REFERENCES people (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT fk_invoice_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table invoice_applied_credits
--
ALTER TABLE invoice_applied_credits
  ADD CONSTRAINT fk_iac_credit_memo FOREIGN KEY (credit_memo_id) REFERENCES credit_memo (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_iac_invoice FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table invoice_applied_discounts
--
ALTER TABLE invoice_applied_discounts
  ADD CONSTRAINT fk_iad_ar_account FOREIGN KEY (ar_account) REFERENCES `accounts` (id),
  ADD CONSTRAINT fk_iad_income_account FOREIGN KEY (income_account) REFERENCES `accounts` (id),
  ADD CONSTRAINT fk_iad_invoice FOREIGN KEY (invoice_id) REFERENCES invoices (id),
  ADD CONSTRAINT fk_iad_transaction FOREIGN KEY (transaction_id) REFERENCES `transactions` (id);

--
-- Constraints for table invoice_payments
--
ALTER TABLE invoice_payments
  ADD CONSTRAINT fk_ip_account FOREIGN KEY (account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_ip_invoice FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_ip_transaction FOREIGN KEY (transaction_id) REFERENCES `transactions` (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_ip_user FOREIGN KEY (user_id) REFERENCES `users` (id) ON UPDATE CASCADE;

--
-- Constraints for table items
--
ALTER TABLE items
  ADD CONSTRAINT fk_items_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table journal
--
ALTER TABLE journal
  ADD CONSTRAINT journal_ibfk_1 FOREIGN KEY (transaction_id) REFERENCES `transactions` (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT journal_ibfk_2 FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table journal_lines
--
ALTER TABLE journal_lines
  ADD CONSTRAINT journal_lines_ibfk_1 FOREIGN KEY (journal_id) REFERENCES journal (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT journal_lines_ibfk_2 FOREIGN KEY (account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT journal_lines_ibfk_3 FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT journal_lines_ibfk_4 FOREIGN KEY (people_id) REFERENCES people (id) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table leases
--
ALTER TABLE leases
  ADD CONSTRAINT fk_leases_buildings FOREIGN KEY (building_id) REFERENCES buildings (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_leases_people FOREIGN KEY (people_id) REFERENCES people (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_leases_units FOREIGN KEY (unit_id) REFERENCES units (id) ON UPDATE CASCADE;

--
-- Constraints for table lease_files
--
ALTER TABLE lease_files
  ADD CONSTRAINT fk_lease_files_lease_id FOREIGN KEY (lease_id) REFERENCES leases (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table people
--
ALTER TABLE people
  ADD CONSTRAINT people_building_id_fk FOREIGN KEY (building_id) REFERENCES buildings (id),
  ADD CONSTRAINT people_type_id_fk FOREIGN KEY (type_id) REFERENCES people_types (id);

--
-- Constraints for table periods
--
ALTER TABLE periods
  ADD CONSTRAINT period_building_id FOREIGN KEY (building_id) REFERENCES buildings (id);

--
-- Constraints for table readings
--
ALTER TABLE readings
  ADD CONSTRAINT fk_readings_item FOREIGN KEY (item_id) REFERENCES items (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_readings_lease FOREIGN KEY (lease_id) REFERENCES leases (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT fk_readings_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON UPDATE CASCADE;

--
-- Constraints for table receipt_items
--
ALTER TABLE receipt_items
  ADD CONSTRAINT fk_sri_item FOREIGN KEY (item_id) REFERENCES items (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_sri_receipt FOREIGN KEY (receipt_id) REFERENCES sales_receipt (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table roles
--
ALTER TABLE roles
  ADD CONSTRAINT fk_roles_owner FOREIGN KEY (owner_user_id) REFERENCES `users` (id) ON DELETE CASCADE;

--
-- Constraints for table role_permissions
--
ALTER TABLE role_permissions
  ADD CONSTRAINT fk_rp_permission FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_rp_role FOREIGN KEY (role_id) REFERENCES `roles` (id) ON DELETE CASCADE;

--
-- Constraints for table sales_receipt
--
ALTER TABLE sales_receipt
  ADD CONSTRAINT fk_sr_account FOREIGN KEY (account_id) REFERENCES `accounts` (id) ON UPDATE CASCADE,
  ADD CONSTRAINT fk_sr_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_sr_people FOREIGN KEY (people_id) REFERENCES people (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT fk_sr_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table splits
--
ALTER TABLE splits
  ADD CONSTRAINT fk_splits_account FOREIGN KEY (account_id) REFERENCES `accounts` (id),
  ADD CONSTRAINT fk_splits_people FOREIGN KEY (people_id) REFERENCES people (id) ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT fk_splits_unit FOREIGN KEY (unit_id) REFERENCES units (id);

--
-- Constraints for table transactions
--
ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_transactions_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table units
--
ALTER TABLE units
  ADD CONSTRAINT unit_building_id_fk FOREIGN KEY (building_id) REFERENCES buildings (id);

--
-- Constraints for table users
--
ALTER TABLE users
  ADD CONSTRAINT parent_user_id FOREIGN KEY (parent_user_id) REFERENCES `users` (id);

--
-- Constraints for table users_building
--
ALTER TABLE users_building
  ADD CONSTRAINT user_building_building_id FOREIGN KEY (building_id) REFERENCES buildings (id),
  ADD CONSTRAINT user_building_user_id FOREIGN KEY (user_id) REFERENCES `users` (id);

--
-- Constraints for table user_building_roles
--
ALTER TABLE user_building_roles
  ADD CONSTRAINT fk_ubr_building FOREIGN KEY (building_id) REFERENCES buildings (id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_ubr_role FOREIGN KEY (role_id) REFERENCES `roles` (id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_ubr_user FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
