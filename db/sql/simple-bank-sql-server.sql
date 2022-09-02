CREATE TABLE [accounts] (
  [id] bigserial PRIMARY KEY,
  [owner] nvarchar(255) NOT NULL,
  [balance] bigint NOT NULL,
  [currency] nvarchar(255) NOT NULL,
  [created_at] timestampz NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [entries] (
  [id] bigserial PRIMARY KEY,
  [account_id] bigint,
  [amount] bigint NOT NULL,
  [created_at] timestampz NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [transfers] (
  [id] bigserial PRIMARY KEY,
  [from_account_id] bigint,
  [to_account_id] bigint,
  [amount] bigint NOT NULL,
  [created_at] timestampz NOT NULL DEFAULT (now())
)
GO

CREATE INDEX [accounts_index_0] ON [accounts] ("owner")
GO

CREATE INDEX [entries_index_1] ON [entries] ("account_id")
GO

CREATE INDEX [transfers_index_2] ON [transfers] ("from_account_id")
GO

CREATE INDEX [transfers_index_3] ON [transfers] ("to_account_id")
GO

CREATE INDEX [transfers_index_4] ON [transfers] ("from_account_id", "to_account_id")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'can be nagative or positive',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entries',
@level2type = N'Column', @level2name = 'amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'must be positive',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'transfers',
@level2type = N'Column', @level2name = 'amount';
GO

ALTER TABLE [entries] ADD FOREIGN KEY ([account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([from_account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([to_account_id]) REFERENCES [accounts] ([id])
GO
