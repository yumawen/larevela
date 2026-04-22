-- Multi-chain Web3 wallet payment starter schema
-- Suggested for services: trade-api, order-rpc, payment-rpc, chain-rpc, ledger-rpc

CREATE DATABASE IF NOT EXISTS `larevela`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `larevela`;

-- 1) Business orders managed by order-rpc
CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(64) NOT NULL COMMENT 'Business order number',
  `biz_type` VARCHAR(32) NOT NULL COMMENT 'Business type, e.g. subscription/product',
  `biz_id` VARCHAR(64) NOT NULL COMMENT 'Business object id',
  `user_id` BIGINT NOT NULL COMMENT 'User id',
  `currency` VARCHAR(16) NOT NULL COMMENT 'Fiat currency, e.g. USD/CNY',
  `amount` DECIMAL(36,18) NOT NULL COMMENT 'Order amount in fiat currency',
  `status` VARCHAR(32) NOT NULL DEFAULT 'created' COMMENT 'created/pending/paid/closed/expired',
  `expired_at` DATETIME NULL COMMENT 'Order expiration time',
  `paid_at` DATETIME NULL COMMENT 'Order paid time',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_biz_type_biz_id` (`biz_type`, `biz_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Business orders';

-- 2) Payment intents managed by payment-rpc
CREATE TABLE IF NOT EXISTS `payment_intents` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_no` VARCHAR(64) NOT NULL COMMENT 'Payment intent number',
  `order_no` VARCHAR(64) NOT NULL COMMENT 'Related business order number',
  `chain_type` VARCHAR(16) NOT NULL DEFAULT 'evm' COMMENT 'evm/solana',
  `network` VARCHAR(32) NOT NULL DEFAULT 'mainnet' COMMENT 'mainnet/testnet/devnet',
  `chain_id` BIGINT NOT NULL COMMENT 'Chain id, e.g. 1/137/42161',
  `pay_mode` VARCHAR(32) NOT NULL DEFAULT 'transfer' COMMENT 'transfer/contract_call',
  `payer_account` VARCHAR(128) NULL COMMENT 'Payer wallet account/address',
  `receiver_account` VARCHAR(128) NOT NULL COMMENT 'Merchant receiver account/address',
  `payer_token_account` VARCHAR(128) NULL COMMENT 'Solana payer token account for SPL token payments',
  `receiver_token_account` VARCHAR(128) NULL COMMENT 'Solana receiver token account for SPL token payments',
  `asset_symbol` VARCHAR(32) NOT NULL COMMENT 'Asset symbol, e.g. ETH/USDT',
  `asset_address` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Token contract or Solana mint address, empty for native coin',
  `amount_expected` DECIMAL(36,18) NOT NULL COMMENT 'Expected payment amount',
  `amount_actual` DECIMAL(36,18) NULL COMMENT 'Actual paid amount',
  `decimals` INT NOT NULL DEFAULT 18 COMMENT 'Token decimals',
  `contract_address` VARCHAR(128) NULL COMMENT 'Checkout contract address',
  `method` VARCHAR(64) NULL COMMENT 'Contract method name',
  `calldata` TEXT NULL COMMENT 'Encoded contract calldata',
  `tx_value` DECIMAL(36,18) NULL COMMENT 'Native value attached to tx',
  `reference_id` VARCHAR(128) NULL COMMENT 'Optional payment reference or memo',
  `quote_expired_at` DATETIME NULL COMMENT 'Quote expiration time',
  `tx_id` VARCHAR(128) NULL COMMENT 'On-chain tx hash or Solana signature',
  `slot` BIGINT NULL COMMENT 'Solana slot number',
  `confirmations` BIGINT NOT NULL DEFAULT 0 COMMENT 'Current confirmation count',
  `confirmation_status` VARCHAR(32) NULL COMMENT 'pending/processed/confirmed/finalized',
  `status` VARCHAR(32) NOT NULL DEFAULT 'created' COMMENT 'created/submitted/confirming/paid/underpaid/overpaid/failed/expired',
  `failure_reason` VARCHAR(255) NULL COMMENT 'Failure or mismatch reason',
  `expired_at` DATETIME NULL COMMENT 'Payment expiration time',
  `paid_at` DATETIME NULL COMMENT 'Payment success time',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_no` (`payment_no`),
  UNIQUE KEY `uk_tx_id` (`tx_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_chain_status` (`chain_type`, `network`, `chain_id`, `status`),
  KEY `idx_payer_account` (`payer_account`),
  KEY `idx_receiver_account` (`receiver_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Payment intents';

-- 3) On-chain transactions and verification results
CREATE TABLE IF NOT EXISTS `payment_transactions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_no` VARCHAR(64) NOT NULL COMMENT 'Related payment number',
  `tx_id` VARCHAR(128) NOT NULL COMMENT 'Unique tx hash or Solana signature',
  `chain_type` VARCHAR(16) NOT NULL DEFAULT 'evm' COMMENT 'evm/solana',
  `network` VARCHAR(32) NOT NULL DEFAULT 'mainnet' COMMENT 'mainnet/testnet/devnet',
  `chain_id` BIGINT NOT NULL COMMENT 'Chain id',
  `from_account` VARCHAR(128) NOT NULL COMMENT 'On-chain sender account/address',
  `to_account` VARCHAR(128) NOT NULL COMMENT 'On-chain receiver account/address',
  `from_token_account` VARCHAR(128) NULL COMMENT 'Solana sender token account',
  `to_token_account` VARCHAR(128) NULL COMMENT 'Solana receiver token account',
  `asset_address` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Token contract or Solana mint address',
  `asset_symbol` VARCHAR(32) NOT NULL COMMENT 'Asset symbol',
  `amount_actual` DECIMAL(36,18) NOT NULL COMMENT 'Actual transferred amount',
  `tx_status` VARCHAR(32) NOT NULL COMMENT 'pending/success/failed/reverted',
  `block_number` BIGINT NULL COMMENT 'Block number containing the tx',
  `slot` BIGINT NULL COMMENT 'Solana slot number',
  `block_hash` VARCHAR(128) NULL COMMENT 'Block hash',
  `gas_used` BIGINT NULL COMMENT 'Gas used from receipt',
  `fee_amount` DECIMAL(36,18) NULL COMMENT 'Gas fee or Solana network fee',
  `confirmations` BIGINT NOT NULL DEFAULT 0 COMMENT 'Observed confirmations',
  `confirmation_status` VARCHAR(32) NULL COMMENT 'pending/processed/confirmed/finalized',
  `reference_id` VARCHAR(128) NULL COMMENT 'Optional memo/reference id',
  `raw_payload` JSON NULL COMMENT 'Raw tx/receipt payload',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chain_tx` (`chain_type`, `network`, `chain_id`, `tx_id`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_from_account` (`from_account`),
  KEY `idx_to_account` (`to_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='On-chain payment transactions';

-- 4) Chain scan state for chain-rpc background sync
CREATE TABLE IF NOT EXISTS `chain_scan` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `chain_type` VARCHAR(16) NOT NULL DEFAULT 'evm' COMMENT 'evm/solana',
  `network` VARCHAR(32) NOT NULL DEFAULT 'mainnet' COMMENT 'mainnet/testnet/devnet',
  `chain_id` BIGINT NOT NULL COMMENT 'Chain id',
  `cursor_type` VARCHAR(32) NOT NULL COMMENT 'block/log/payment_confirm',
  `last_scanned_block` BIGINT NOT NULL DEFAULT 0 COMMENT 'Last scanned block number',
  `last_scanned_slot` BIGINT NOT NULL DEFAULT 0 COMMENT 'Last scanned Solana slot',
  `last_scanned_tx_id` VARCHAR(128) NULL COMMENT 'Last processed tx id if needed',
  `remark` VARCHAR(255) NULL COMMENT 'Additional notes',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chain_cursor_type` (`chain_type`, `network`, `chain_id`, `cursor_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Chain scan cursors';

-- 5) Ledger entries managed by ledger-rpc
CREATE TABLE IF NOT EXISTS `ledger_entries` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `entry_no` VARCHAR(64) NOT NULL COMMENT 'Ledger entry number',
  `payment_no` VARCHAR(64) NOT NULL COMMENT 'Related payment number',
  `order_no` VARCHAR(64) NOT NULL COMMENT 'Related order number',
  `user_id` BIGINT NOT NULL COMMENT 'User id',
  `chain_type` VARCHAR(16) NOT NULL DEFAULT 'evm' COMMENT 'evm/solana',
  `network` VARCHAR(32) NOT NULL DEFAULT 'mainnet' COMMENT 'mainnet/testnet/devnet',
  `chain_id` BIGINT NOT NULL COMMENT 'Chain id',
  `entry_type` VARCHAR(32) NOT NULL COMMENT 'payment/refund/adjustment',
  `asset_symbol` VARCHAR(32) NOT NULL COMMENT 'Asset symbol',
  `asset_address` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Token contract or Solana mint address',
  `amount` DECIMAL(36,18) NOT NULL COMMENT 'Ledger amount',
  `direction` VARCHAR(16) NOT NULL COMMENT 'debit/credit',
  `status` VARCHAR(32) NOT NULL DEFAULT 'posted' COMMENT 'pending/posted/reversed',
  `remark` VARCHAR(255) NULL COMMENT 'Business remark',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_entry_no` (`entry_no`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_chain_asset` (`chain_type`, `network`, `chain_id`, `asset_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Ledger entries';

-- 6) Idempotency table for repeated submissions and async retries
CREATE TABLE IF NOT EXISTS `idempotency_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `idem_key` VARCHAR(128) NOT NULL COMMENT 'Idempotency key',
  `biz_type` VARCHAR(32) NOT NULL COMMENT 'submit_tx/create_intent/ledger_post/etc',
  `biz_no` VARCHAR(64) NOT NULL COMMENT 'Business number',
  `status` VARCHAR(32) NOT NULL COMMENT 'processing/success/failed',
  `response_snapshot` JSON NULL COMMENT 'Optional response snapshot for replay',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_idem_key` (`idem_key`),
  KEY `idx_biz_type_biz_no` (`biz_type`, `biz_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Idempotency records';

-- 7) Aggregated payment read model for frontend status display
CREATE TABLE IF NOT EXISTS `payment_view` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_no` VARCHAR(64) NOT NULL COMMENT 'Payment number',
  `order_no` VARCHAR(64) NOT NULL COMMENT 'Related order number',
  `tx_id` VARCHAR(128) NULL COMMENT 'Tx hash or Solana signature',
  `chain_type` VARCHAR(16) NOT NULL DEFAULT 'solana' COMMENT 'evm/solana',
  `network` VARCHAR(32) NOT NULL DEFAULT 'devnet' COMMENT 'mainnet/testnet/devnet',
  `chain_id` BIGINT NOT NULL COMMENT 'Chain id',
  `payer_account` VARCHAR(128) NULL COMMENT 'Payer wallet account/address',
  `receiver_account` VARCHAR(128) NULL COMMENT 'Receiver wallet account/address',
  `asset_symbol` VARCHAR(32) NOT NULL COMMENT 'Asset symbol',
  `amount_expected` DECIMAL(36,18) NOT NULL COMMENT 'Expected payment amount',
  `amount_actual` DECIMAL(36,18) NULL COMMENT 'Actual paid amount',
  `status` VARCHAR(32) NOT NULL DEFAULT 'created' COMMENT 'created/submitted/confirming/paid/failed',
  `confirmation_status` VARCHAR(32) NULL COMMENT 'pending/processed/confirmed/finalized',
  `confirmations` BIGINT NOT NULL DEFAULT 0 COMMENT 'Confirmation count',
  `last_scanned_block` BIGINT NOT NULL DEFAULT 0 COMMENT 'Last scanned block',
  `last_scanned_slot` BIGINT NOT NULL DEFAULT 0 COMMENT 'Last scanned slot',
  `failure_reason` VARCHAR(255) NULL COMMENT 'Failure reason',
  `updated_source` VARCHAR(32) NOT NULL DEFAULT 'payment' COMMENT 'payment/confirm/sync',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_no` (`payment_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Aggregated payment status view';
