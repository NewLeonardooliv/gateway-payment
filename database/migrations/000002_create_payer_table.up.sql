CREATE TABLE IF NOT EXISTS payers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    name TEXT NOT NULL,
    tax_id TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    address TEXT NOT NULL,
    number TEXT NOT NULL,
    district TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    zip_code TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL
);

CREATE INDEX idx_payer_invoice_id ON payers(invoice_id);
CREATE INDEX idx_payer_created_at ON payers(created_at);