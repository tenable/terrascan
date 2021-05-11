package functions

import (
	"testing"

	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLookUp(t *testing.T) {
	vars := map[string]interface{}{
		"sslEnabled":       true,
		"dbName":           "database-name",
		"sqlAdminPassword": "sql-Admin-Password",
	}

	params := map[string]interface{}{
		"retentionDays":         7,
		"sqlAdministratorLogin": "sql-admin",
	}

	const resourceType = "Microsoft.Sql/servers"
	resourceConfigID := types.ResourceTypes[resourceType] + ".sql-server"
	ResourceIDs[resourceType] = resourceConfigID

	tests := []struct {
		name     string
		key      string
		expected interface{}
	}{
		{
			name:     "variables-string-value",
			key:      "dbName",
			expected: "database-name",
		},
		{
			name:     "variables-bool-value",
			key:      "sslEnabled",
			expected: true,
		},
		{
			name:     "parameters-string-value",
			key:      "retentionDays",
			expected: 7,
		},
		{
			name:     "unknown-key",
			key:      "unknown-key",
			expected: "unknown-key",
		},
		{
			name:     "parameters",
			key:      "[parameters('sqlAdministratorLogin')]",
			expected: "sql-admin",
		},
		{
			name:     "parameters",
			key:      "[parameters'sqlAdministratorLogin')]",
			expected: "",
		},
		{
			name:     "variables",
			key:      "[variables('dbName')]",
			expected: "database-name",
		},
		{
			name:     "variables-invalid",
			key:      "[variables'dbName')]",
			expected: "",
		},
		{
			name:     "concat-variables",
			key:      "[concat('NetworkWatcher_', variables('dbName'))]",
			expected: "NetworkWatcher_database-name",
		},
		{
			name:     "concat-variables-parameters",
			key:      "[concat('NetworkWatcher_', variables('dbName'), '_', parameters('sqlAdministratorLogin'))]",
			expected: "NetworkWatcher_database-name_sql-admin",
		},
		{
			name:     "concat-invalid",
			key:      "[concat'NetworkWatcher_', variables('dbName'))]",
			expected: "",
		},
		{
			name:     "toLower",
			key:      "toLower(variables('sqlAdminPassword'))",
			expected: "sql-admin-password",
		},
		{
			name:     "toLower-invalid",
			key:      "toLowervariables('sqlAdminPassword'))",
			expected: "",
		},
		{
			name:     "resourceId",
			key:      "[resourceId('Microsoft.Sql/servers', variables('sqlServerName'))]",
			expected: resourceConfigID,
		},
		{
			name:     "resourceId-invalid",
			key:      "[resourceId'Microsoft.Sql/servers', variables('sqlServerName'))]",
			expected: "",
		},
		{
			name:     "resourceId-not-found",
			key:      "[resourceId('Microsoft.KeyVault/vaults', parameters('keyVaultName'))]",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := LookUp(vars, params, test.key)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestUniqueString(t *testing.T) {
	const str = "uniqueString(resourceGroup().id)"
	res := LookUp(nil, nil, str).(string)

	_, err := uuid.Parse(res)
	assert.NoError(t, err)
}
