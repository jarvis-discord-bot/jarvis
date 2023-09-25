CREATE TABLE IF NOT EXISTS chat_commands
(
  guild_id      VARCHAR(32),
  command_key   VARCHAR(64),
  command_value TEXT,
  PRIMARY KEY (guild_id, command_key)
);
