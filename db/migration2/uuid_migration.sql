ALTER TABLE contact_infos ADD COLUMN uuid UUID;

-- UPDATE contact_infos SET uuid = gen_random_uuid() WHERE uuid IS NULL;

UPDATE contact_infos SET uuid = NULL;

select uuid from contact_infos where uuid is null;

ALTER TABLE contact_infos ALTER COLUMN uuid SET NOT NULL;

ALTER TABLE contact_infos ADD CONSTRAINT contact_infos_uuid_unique UNIQUE (uuid);

