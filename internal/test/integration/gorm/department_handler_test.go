// Package integration_test ...
package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- helpers ---

func toJSON(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func decodeJSON(t *testing.T, body *bytes.Buffer, dst any) {
	t.Helper()
	require.NoError(t, json.NewDecoder(body).Decode(dst))
}

func createDepartment(t *testing.T, name string, parentID *uint) map[string]any {
	t.Helper()
	body := map[string]any{"name": name}
	if parentID != nil {
		body["parent_id"] = *parentID
	}
	req, _ := http.NewRequest(http.MethodPost, "/departments/", toJSON(body))
	req.Header.Set("Content-Type", "application/json")
	rec := app.do(req)
	require.Equal(t, http.StatusCreated, rec.Code)
	var result map[string]any
	decodeJSON(t, rec.Body, &result)
	return result
}

func idFrom(t *testing.T, m map[string]any) uint {
	t.Helper()
	return uint(m["id"].(float64))
}

// --- Department ---

func TestIntegration_CreateDepartment_Success(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/departments/", toJSON(map[string]any{
		"name": "Engineering",
	}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp map[string]any
	decodeJSON(t, rec.Body, &resp)
	assert.Equal(t, "Engineering", resp["name"])
	assert.NotNil(t, resp["id"])
}

func TestIntegration_CreateDepartment_EmptyName(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/departments/", toJSON(map[string]any{
		"name": "   ",
	}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIntegration_CreateDepartment_DuplicateName(t *testing.T) {
	createDepartment(t, "DuplicateDept", nil)

	req, _ := http.NewRequest(http.MethodPost, "/departments/", toJSON(map[string]any{
		"name": "DuplicateDept",
	}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestIntegration_GetDepartment_WithChildren(t *testing.T) {
	parent := createDepartment(t, "ParentDept", nil)
	parentID := idFrom(t, parent)
	createDepartment(t, "ChildDept", &parentID)

	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("/departments/%d?depth=2&include_employees=true", parentID), nil)

	rec := app.do(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	decodeJSON(t, rec.Body, &resp)
	assert.Equal(t, "ParentDept", resp["name"])
	children := resp["children"].([]any)
	assert.Len(t, children, 1)
	assert.Equal(t, "ChildDept", children[0].(map[string]any)["name"])
}

func TestIntegration_GetDepartment_NotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/departments/999999", nil)
	rec := app.do(req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIntegration_UpdateDepartment_Rename(t *testing.T) {
	dept := createDepartment(t, "OldName", nil)
	id := idFrom(t, dept)

	req, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("/departments/%d", id),
		toJSON(map[string]any{"name": "NewName"}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	decodeJSON(t, rec.Body, &resp)
	assert.Equal(t, "NewName", resp["name"])
}

func TestIntegration_UpdateDepartment_CyclicDependency(t *testing.T) {
	parent := createDepartment(t, "CycleParent", nil)
	parentID := idFrom(t, parent)
	child := createDepartment(t, "CycleChild", &parentID)
	childID := idFrom(t, child)

	req, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("/departments/%d", parentID),
		toJSON(map[string]any{"parent_id": childID}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestIntegration_DeleteDepartment_Cascade(t *testing.T) {
	dept := createDepartment(t, "ToDeleteCascade", nil)
	id := idFrom(t, dept)

	req, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("/departments/%d?mode=cascade", id), nil)
	rec := app.do(req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	getReq, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("/departments/%d", id), nil)
	getRec := app.do(getReq)
	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestIntegration_DeleteDepartment_Reassign(t *testing.T) {
	target := createDepartment(t, "ReassignTarget", nil)
	targetID := idFrom(t, target)

	toDelete := createDepartment(t, "ToDeleteReassign", nil)
	toDeleteID := idFrom(t, toDelete)

	empReq, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", toDeleteID),
		toJSON(map[string]any{
			"full_name": "Test Employee",
			"position":  "Developer",
		}))
	empReq.Header.Set("Content-Type", "application/json")
	empRec := app.do(empReq)
	require.Equal(t, http.StatusCreated, empRec.Code)

	req, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("/departments/%d?mode=reassign&reassign_to_department_id=%d", toDeleteID, targetID), nil)
	rec := app.do(req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	getReq, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("/departments/%d?include_employees=true", targetID), nil)
	getRec := app.do(getReq)
	assert.Equal(t, http.StatusOK, getRec.Code)
	var resp map[string]any
	decodeJSON(t, getRec.Body, &resp)
	employees := resp["employees"].([]any)
	assert.Len(t, employees, 1)
	assert.Equal(t, "Test Employee", employees[0].(map[string]any)["full_name"])
}

// --- Employee ---

func TestIntegration_CreateEmployee_Success(t *testing.T) {
	dept := createDepartment(t, "EmpDept", nil)
	deptID := idFrom(t, dept)

	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", deptID),
		toJSON(map[string]any{
			"full_name": "Ivan Petrov",
			"position":  "Senior Developer",
		}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp map[string]any
	decodeJSON(t, rec.Body, &resp)
	assert.Equal(t, "Ivan Petrov", resp["full_name"])
	assert.Equal(t, "Senior Developer", resp["position"])
}

func TestIntegration_CreateEmployee_DepartmentNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost,
		"/departments/999999/employees/",
		toJSON(map[string]any{
			"full_name": "Ghost",
			"position":  "Nobody",
		}))
	req.Header.Set("Content-Type", "application/json")

	rec := app.do(req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
