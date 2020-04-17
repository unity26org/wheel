package templatecommon

var MigrateContent = `var m{{ .MigrationMetadata.Version }} migrate.Version{{ .MigrationMetadata.Version }}
root = pushToTree(root, "{{ .MigrationMetadata.Version }}", m{{ .MigrationMetadata.Version }}.{{ .MigrationMetadata.Name }})`
