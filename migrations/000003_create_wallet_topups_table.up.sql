CREATE TABLE wallet_topups (
    id BIGINT NOT NULL AUTO_INCREMENT,
    wallet_id BIGINT NOT NULL,
    order_id VARCHAR(64) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    snap_token VARCHAR(255) NOT NULL DEFAULT '',
    redirect_url VARCHAR(1024) NOT NULL DEFAULT '',
    payment_type VARCHAR(50) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY wallet_topups_order_id_unique (order_id),
    KEY wallet_topups_wallet_id_index (wallet_id),
    CONSTRAINT wallet_topups_wallet_id_foreign
        FOREIGN KEY (wallet_id) REFERENCES wallets (id)
        ON DELETE CASCADE
);
