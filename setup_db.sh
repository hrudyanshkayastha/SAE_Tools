su - postgres -c "psql -c \\"CREATE USER sae WITH PASSWORD 'changeme_dev_only';\\""
su - postgres -c "psql -c \\"CREATE DATABASE sae_db OWNER sae;\\""
