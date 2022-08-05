create table accounts (
  account_id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  data jsonb
);

create table documents (
  document_id uuid DEFAULT gen_random_uuid (),
  account_id uuid REFERENCES accounts ON DELETE CASCADE,
  data jsonb,
  PRIMARY KEY(account_id, document_id)
);