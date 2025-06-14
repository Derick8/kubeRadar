package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kubeRadar/pkg/models"

	"github.com/xuri/excelize/v2"

	// Import image format packages for excelize to handle different image types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Report represents an Excel report generator
type Report struct {
	filePath      string
	excel         *excelize.File
	headerStyle   int
	contentStyle  int
	sectionStyle  int
	altRowStyle   int
	wrapTextStyle int
	titleStyle    int
	criticalStyle int
	warningStyle  int
	moderateStyle int
	goodStyle     int

	// Chart styles
	chartTitleStyleID  int
	chartLegendStyleID int
	chartAxisStyleID   int
	chartSeriesColors  []string
}

func NewReport(filePath string) (*Report, error) {
	r := &Report{
		excel:    excelize.NewFile(),
		filePath: filePath,
	}

	// Initialize styles
	err := r.setupStyles()
	if err != nil {
		return nil, fmt.Errorf("failed to setup styles: %v", err)
	}

	return r, nil
}

// setupStyles creates and sets up the styles used in the report
func (r *Report) setupStyles() error {
	var err error

	// Initialize chart styles
	r.chartSeriesColors = []string{
		"FF4B55",
		"FF9800",
		"FFD700",
		"4CAF50",
		"2196F3",
	}

	// Title style - Large bold text
	r.titleStyle, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create title style: %v", err)
	}

	// Header style - Blue background with white text
	r.headerStyle, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %v", err)
	}

	// Content style - Basic bordered cells
	r.contentStyle, err = r.excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create content style: %v", err)
	}

	// Section style - Light green background
	r.sectionStyle, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"E2EFDA"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create section style: %v", err)
	}

	// Alternating row style - Light gray background
	r.altRowStyle, err = r.excel.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F5F5F5"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create alternating row style: %v", err)
	}

	// Wrap text style - For cells with long content
	r.wrapTextStyle, err = r.excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create wrap text style: %v", err)
	}

	r.criticalStyle, err = r.excel.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFCCCC"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create critical style: %v", err)
	}

	// Warning style - Yellow background for warnings
	r.warningStyle, err = r.excel.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFFCC"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create warning style: %v", err)
	}

	// Moderate style - Light green background for moderate issues
	r.moderateStyle, err = r.excel.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"CCFFCC"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create moderate style: %v", err)
	}

	// Good style - Green background for good status
	r.goodStyle, err = r.excel.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"D9EAD3"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create good style: %v", err)
	}

	// Initialize chart styles
	r.chartSeriesColors = []string{
		"FF4B55",
		"FF9800",
		"FFD700",
		"4CAF50",
		"2196F3",
	}

	r.chartTitleStyleID, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "000000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create chart title style: %v", err)
	}

	r.chartLegendStyleID, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "000000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create chart legend style: %v", err)
	}
	r.chartAxisStyleID, err = r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create chart axis style: %v", err)
	}

	return nil
}

// formatLabels formats a map of labels into a string
func (r *Report) formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}

	var result []string
	for k, v := range labels {
		result = append(result, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(result)
	return strings.Join(result, "\n")
}

// formatCells applies styles to a range of cells

// autoFitColumns automatically adjusts column widths in a sheet
func (r *Report) autoFitColumns(sheet string) error {
	cols, err := r.excel.GetCols(sheet)
	if err != nil {
		return err
	}

	for idx, col := range cols {
		maxLen := 0
		for _, cell := range col {
			if len(cell) > maxLen {
				maxLen = len(cell)
			}
		}

		colName, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			continue
		}

		// Set column width with some padding
		width := float64(maxLen) + 2
		if width > 100 {
			width = 100 // Set maximum width
		} else if width < 10 {
			width = 10 // Set minimum width
		}
		r.excel.SetColWidth(sheet, colName, colName, width)
	}

	return nil
}

func (r *Report) Generate(data *models.AssessmentData) error {
	// Setup styles first
	if err := r.setupStyles(); err != nil {
		return err
	}
	// Create all sheets
	sheets := []string{
		"Contents",
		"Dashboard",
		"Nodes",
		"Namespaces",
		"Pods",
		"Deployments",
		"StatefulSets",
		"DaemonSets",
		"Services",
		"Network Policies",
		"Ingresses",
		"Secrets",
		"Service Accounts",
		"Roles",
		"Role Bindings",
		"Cluster Roles",
		"Cluster Role Bindings",
	}

	// Initialize sheets
	for i, sheet := range sheets {
		if i == 0 {
			r.excel.SetSheetName("Sheet1", sheet)
		} else {
			r.excel.NewSheet(sheet)
		}
	}

	// Generate Contents page first
	if err := r.generateTableOfContents(sheets); err != nil {
		return fmt.Errorf("failed to generate table of contents: %v", err)
	}

	if err := r.generateDashboard(data); err != nil {
		return fmt.Errorf("failed to generate dashboard: %v", err)
	}
	if err := r.generateNodes(data.ClusterInfo.Nodes); err != nil {
		return fmt.Errorf("failed to generate nodes: %v", err)
	}
	if err := r.generateNamespaces(data.ClusterInfo.Namespaces); err != nil {
		return fmt.Errorf("failed to generate namespaces: %v", err)
	}
	if err := r.generatePods(data.Workloads); err != nil {
		return fmt.Errorf("failed to generate pods: %v", err)
	}
	if err := r.generateDeployments(data.Workloads.Deployments); err != nil {
		return fmt.Errorf("failed to generate deployments: %v", err)
	}
	if err := r.generateStatefulSets(data.Workloads.StatefulSets); err != nil {
		return fmt.Errorf("failed to generate stateful sets: %v", err)
	}
	if err := r.generateDaemonSets(data.Workloads.DaemonSets); err != nil {
		return fmt.Errorf("failed to generate daemon sets: %v", err)
	}
	if err := r.generateServices(data.Network.Services); err != nil {
		return fmt.Errorf("failed to generate services: %v", err)
	}
	if err := r.generateNetworkPolicies(data.Network.NetworkPolicies); err != nil {
		return fmt.Errorf("failed to generate network policies: %v", err)
	}
	if err := r.generateIngresses(data.Network.Ingresses); err != nil {
		return fmt.Errorf("failed to generate ingresses: %v", err)
	}
	if err := r.generateSecrets(data.Secrets.Secrets); err != nil {
		return fmt.Errorf("failed to generate secrets: %v", err)
	}
	if err := r.generateServiceAccounts(data.RBAC.ServiceAccounts); err != nil {
		return fmt.Errorf("failed to generate service accounts: %v", err)
	}
	// Only keep the following calls for RBAC:
	if err := r.generateRoles(data.RBAC); err != nil {
		return fmt.Errorf("failed to generate roles: %v", err)
	}
	if err := r.generateRoleBindings(data.RBAC); err != nil {
		return fmt.Errorf("failed to generate role bindings: %v", err)
	}
	if err := r.generateClusterRoles(data.RBAC); err != nil {
		return fmt.Errorf("failed to generate cluster roles: %v", err)
	}
	if err := r.generateClusterRoleBindings(data.RBAC); err != nil {
		return fmt.Errorf("failed to generate cluster role bindings: %v", err)
	}

	// Auto-fit columns in all sheets
	for _, sheet := range sheets {
		if err := r.autoFitColumns(sheet); err != nil {
			return fmt.Errorf("failed to auto-fit columns in %s: %v", sheet, err)
		}
	}

	// Save the file
	if err := r.excel.SaveAs(r.filePath); err != nil {
		return fmt.Errorf("failed to save excel file: %v", err)
	}

	return nil
}

func (r *Report) generateTableOfContents(_ []string) error {
	sheet := "Contents"
	// Insert logo image at the top (cell A1) using file-based approach
	logoPath := filepath.Join("pkg", "excel", "logo.png")
	if err := r.excel.AddPicture(sheet, "A1", logoPath, &excelize.GraphicOptions{
		OffsetX: 0, OffsetY: 0, ScaleX: 0.28, ScaleY: 0.79,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "[kubeRadar] Warning: could not insert logo from %s: %v\n", logoPath, err)
	}

	// Title
	r.excel.SetCellValue(sheet, "A8", "Kubernetes Reconnaissance Report")
	r.excel.MergeCell(sheet, "A8", "C8")
	r.excel.SetCellStyle(sheet, "A8", "C8", r.titleStyle)

	// Section: Table of Contents
	r.excel.SetCellValue(sheet, "A10", "Contents")
	r.excel.SetCellStyle(sheet, "A10", "A10", r.sectionStyle)

	// List of sections with hyperlinks
	toc := []struct {
		name   string
		target string
	}{
		{"Dashboard", "Dashboard"},
		{"Nodes", "Nodes"},
		{"Namespaces", "Namespaces"},
		{"Pods", "Pods"},
		{"Deployments", "Deployments"},
		{"StatefulSets", "StatefulSets"},
		{"DaemonSets", "DaemonSets"},
		{"Services", "Services"},
		{"Network Policies", "Network Policies"},
		{"Ingresses", "Ingresses"},
		{"Secrets", "Secrets"},
		{"Service Accounts", "Service Accounts"},
		{"Roles", "Roles"},
		{"Role Bindings", "Role Bindings"},
		{"Cluster Roles", "Cluster Roles"},
		{"Cluster Role Bindings", "Cluster Role Bindings"},
	}
	for i, entry := range toc {
		cell := fmt.Sprintf("A%d", 12+i)
		r.excel.SetCellValue(sheet, cell, entry.name)
		r.excel.SetCellHyperLink(sheet, cell, fmt.Sprintf("#'%s'!A1", entry.target), "Location")
		r.excel.SetCellStyle(sheet, cell, cell, r.contentStyle)
	}
	r.autoFitColumns(sheet)
	return nil
}

// Roles pane
func (r *Report) generateRoles(rbac models.RBACAssessment) error {
	sheet := "Roles"
	headers := []string{"Name", "Namespace", "Created At", "Rules"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, role := range rbac.Roles {
		values := []interface{}{
			role.Name,
			role.Namespace,
			role.CreatedAt,
			FormatRules(role.Rules),
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}
func FormatRules(rules []models.PolicyRule) string {
	var ruleStrings []string
	for _, rule := range rules {
		ruleStr := fmt.Sprintf("API Groups: [%s]\nResources: [%s]\nVerbs: [%s]",
			strings.Join(rule.APIGroups, ", "),
			strings.Join(rule.Resources, ", "),
			strings.Join(rule.Verbs, ", "))
		if len(rule.ResourceNames) > 0 {
			ruleStr += fmt.Sprintf("\nResource Names: [%s]", strings.Join(rule.ResourceNames, ", "))
		}
		ruleStrings = append(ruleStrings, ruleStr)
	}
	return strings.Join(ruleStrings, "\n---\n")
}

// Role Bindings pane
func (r *Report) generateRoleBindings(rbac models.RBACAssessment) error {
	sheet := "Role Bindings"
	headers := []string{"Name", "Namespace", "Role Ref", "Subjects", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, binding := range rbac.RoleBindings {
		subjects := make([]string, 0)
		for _, subject := range binding.Subjects {
			subjects = append(subjects, fmt.Sprintf("%s/%s (%s)", subject.Namespace, subject.Name, subject.Kind))
		}
		values := []interface{}{
			binding.Name,
			binding.Namespace,
			binding.RoleRef,
			strings.Join(subjects, ", "),
			binding.CreatedAt,
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}

// Cluster Roles pane
func (r *Report) generateClusterRoles(rbac models.RBACAssessment) error {
	sheet := "Cluster Roles"
	headers := []string{"Name", "Created At", "Rules"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, role := range rbac.ClusterRoles {
		values := []interface{}{
			role.Name,
			role.CreatedAt,
			FormatRules(role.Rules),
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}

// Cluster Role Bindings pane
func (r *Report) generateClusterRoleBindings(rbac models.RBACAssessment) error {
	sheet := "Cluster Role Bindings"
	headers := []string{"Name", "Role Ref", "Subjects", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, binding := range rbac.ClusterRoleBindings {
		subjects := make([]string, 0)
		for _, subject := range binding.Subjects {
			subjects = append(subjects, fmt.Sprintf("%s/%s (%s)", subject.Namespace, subject.Name, subject.Kind))
		}
		values := []interface{}{
			binding.Name,
			binding.RoleRef,
			strings.Join(subjects, ", "),
			binding.CreatedAt,
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}

// Basic report generation methods for each resource type
func (r *Report) generateDeployments(deployments []models.DeploymentInfo) error {
	sheet := "Deployments"
	headers := []string{"Name", "Namespace", "Replicas", "Update Strategy", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, deploy := range deployments {
		row := i + 2
		values := []interface{}{
			deploy.Name,
			deploy.Namespace,
			deploy.Replicas,
			deploy.UpdateStrategy,
			r.formatLabels(deploy.Labels),
			deploy.CreatedAt,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateStatefulSets(statefulSets []models.StatefulSetInfo) error {
	sheet := "StatefulSets"
	headers := []string{"Name", "Namespace", "Replicas", "Update Strategy", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, sts := range statefulSets {
		row := i + 2
		values := []interface{}{
			sts.Name,
			sts.Namespace,
			sts.Replicas,
			sts.UpdateStrategy,
			r.formatLabels(sts.Labels),
			sts.CreatedAt,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateDaemonSets(daemonSets []models.DaemonSetInfo) error {
	sheet := "DaemonSets"
	headers := []string{"Name", "Namespace", "Update Strategy", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, ds := range daemonSets {
		row := i + 2
		values := []interface{}{
			ds.Name,
			ds.Namespace,
			ds.UpdateStrategy,
			r.formatLabels(ds.Labels),
			ds.CreatedAt,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateServices(services []models.ServiceInfo) error {
	sheet := "Services"
	headers := []string{"Name", "Namespace", "Type", "Cluster IP", "External IP", "Ports", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, svc := range services {
		row := i + 2
		values := []interface{}{
			svc.Name,
			svc.Namespace,
			svc.Type,
			svc.ClusterIP,
			strings.Join(svc.ExternalIPs, ", "),
			r.formatPorts(svc.Ports),
			r.formatLabels(svc.Labels),
			svc.CreatedAt, // <-- add missing comma here
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

// Network Policies pane
func (r *Report) generateNetworkPolicies(networkPolicies []models.NetworkPolicyInfo) error {
	sheet := "Network Policies"
	headers := []string{"Name", "Namespace", "Pod Selector", "Policy Types", "Created At", "Labels"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, policy := range networkPolicies {
		values := []interface{}{
			policy.Name,
			policy.Namespace,
			policy.PodSelector,
			strings.Join(policy.PolicyTypes, ", "),
			policy.CreatedAt,
			r.formatLabels(policy.Labels),
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}

// Service Accounts pane
func (r *Report) generateServiceAccounts(serviceAccounts []models.ServiceAccountInfo) error {
	sheet := "Service Accounts" // Ensure this matches the sheet created in Generate
	// If the sheet does not exist, create it (defensive)
	if idx, _ := r.excel.GetSheetIndex(sheet); idx == -1 {
		r.excel.NewSheet(sheet)
	}
	headers := []string{"Name", "Namespace", "Secrets", "Image Pull Secrets", "Created At", "Labels"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	r.excel.AutoFilter(sheet, fmt.Sprintf("A1:%s1", endCol), nil)
	row := 2
	for _, sa := range serviceAccounts {
		values := []interface{}{
			sa.Name,
			sa.Namespace,
			strings.Join(sa.Secrets, ", "),
			strings.Join(sa.ImagePullSecrets, ", "),
			sa.CreatedAt,
			r.formatLabels(sa.Labels),
		}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
		row++
	}
	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateIngresses(ingresses []models.IngressInfo) error {
	sheet := "Ingresses"
	headers := []string{"Name", "Namespace", "Rules", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, ing := range ingresses {
		row := i + 2
		values := []interface{}{
			ing.Name,
			ing.Namespace,
			r.formatIngressRules(ing.Rules),
			r.formatLabels(ing.Labels),
			ing.CreatedAt,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateSecrets(secrets []models.SecretInfo) error {
	sheet := "Secrets"
	headers := []string{"Name", "Namespace", "Type", "Labels", "Created At"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	for i, secret := range secrets {
		row := i + 2
		values := []interface{}{
			secret.Name,
			secret.Namespace,
			secret.Type,
			r.formatLabels(secret.Labels),
			secret.CreatedAt,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

// Helper methods for formatting
func (r *Report) formatPorts(ports []models.ServicePort) string {
	var portStrings []string
	for _, port := range ports {
		portStr := fmt.Sprintf("%d", port.Port)
		if port.TargetPort != 0 {
			portStr += fmt.Sprintf("→%d", port.TargetPort)
		}
		if port.Protocol != "" {
			portStr += fmt.Sprintf("/%s", port.Protocol)
		}
		portStrings = append(portStrings, portStr)
	}
	return strings.Join(portStrings, "\n")
}

func (r *Report) formatIngressRules(rules []models.IngressRule) string {
	var ruleStrings []string
	for _, rule := range rules {
		ruleStr := rule.Host + " → "
		var paths []string
		for _, path := range rule.Paths {
			paths = append(paths, fmt.Sprintf("%s:%d%s", path.ServiceName, path.ServicePort, path.Path))
		}
		ruleStr += strings.Join(paths, ", ")
		ruleStrings = append(ruleStrings, ruleStr)
	}
	return strings.Join(ruleStrings, "\n")
}

func (r *Report) generateDashboard(data *models.AssessmentData) error {
	sheet := "Dashboard"
	r.excel.SetCellValue(sheet, "A1", "Kubernetes Cluster Configuration Overview")
	r.excel.MergeCell(sheet, "A1", "C1")
	r.excel.SetCellStyle(sheet, "A1", "C1", r.titleStyle)

	// Add cluster information section
	r.excel.SetCellValue(sheet, "A3", "Cluster Overview")
	r.excel.MergeCell(sheet, "A3", "C3")
	r.excel.SetCellStyle(sheet, "A3", "C3", r.sectionStyle)

	// Add key metrics
	metrics := []struct {
		label string
		value interface{}
	}{
		{"Kubernetes Version", data.ClusterInfo.Version},
		{"Total Nodes", data.ClusterInfo.NodeCount},
		{"Total Namespaces", len(data.ClusterInfo.Namespaces)},
		{"Total Pods", len(data.Workloads.Pods)},
		{"Total Deployments", len(data.Workloads.Deployments)},
		{"Total StatefulSets", len(data.Workloads.StatefulSets)},
		{"Total DaemonSets", len(data.Workloads.DaemonSets)},
		{"Total Services", len(data.Network.Services)},
		{"Total Network Policies", len(data.Network.NetworkPolicies)},
		{"Total Ingresses", len(data.Network.Ingresses)},
		{"Total Secrets", len(data.Secrets.Secrets)},
		{"Total Roles", len(data.RBAC.Roles)},
		{"Total ClusterRoles", len(data.RBAC.ClusterRoles)},
		{"Total RoleBindings", len(data.RBAC.RoleBindings)},
		{"Total ClusterRoleBindings", len(data.RBAC.ClusterRoleBindings)},
		{"Total ServiceAccounts", len(data.RBAC.ServiceAccounts)},
	}

	// Add metrics to dashboard
	for i, metric := range metrics {
		row := i + 4
		r.excel.SetCellValue(sheet, fmt.Sprintf("A%d", row), metric.label)
		r.excel.SetCellValue(sheet, fmt.Sprintf("B%d", row), metric.value)
		// Apply alternating row colors
		style := r.contentStyle
		if i%2 == 1 {
			style = r.altRowStyle
		}
		r.excel.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), style)
	}

	// --- RBAC Summary Table and Chart ---
	rbacTableStart := 4
	rbacLabels := []string{"Roles", "ClusterRoles", "RoleBindings", "ClusterRoleBindings", "ServiceAccounts"}
	rbacCounts := []int{
		len(data.RBAC.Roles),
		len(data.RBAC.ClusterRoles),
		len(data.RBAC.RoleBindings),
		len(data.RBAC.ClusterRoleBindings),
		len(data.RBAC.ServiceAccounts),
	}
	r.excel.SetCellValue(sheet, "E3", "RBAC Summary")
	r.excel.MergeCell(sheet, "E3", "F3")
	r.excel.SetCellStyle(sheet, "E3", "F3", r.sectionStyle)
	r.excel.SetCellValue(sheet, "E4", "Type")
	r.excel.SetCellValue(sheet, "F4", "Count")
	r.excel.SetCellStyle(sheet, "E4", "E4", r.headerStyle)
	r.excel.SetCellStyle(sheet, "F4", "F4", r.headerStyle)
	for i, label := range rbacLabels {
		row := rbacTableStart + 1 + i
		r.excel.SetCellValue(sheet, fmt.Sprintf("E%d", row), label)
		r.excel.SetCellValue(sheet, fmt.Sprintf("F%d", row), rbacCounts[i])
		style := r.contentStyle
		if i%2 == 1 {
			style = r.altRowStyle
		}
		r.excel.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), style)
	}
	// Insert RBAC chart (beside table)
	r.excel.AddChart(sheet, "H4", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{{
			Name:       "RBAC Objects",
			Categories: fmt.Sprintf("%s!$E$5:$E$9", sheet),
			Values:     fmt.Sprintf("%s!$F$5:$F$9", sheet),
		}},
		Title:  []excelize.RichTextRun{{Text: "RBAC Objects Distribution"}},
		Legend: excelize.ChartLegend{Position: "top"},
	})

	// --- Pod Security Summary Table and Chart ---
	podTableStart := 12
	podTotal := len(data.Workloads.Pods)
	podPrivileged := 0
	podHostNetwork := 0
	podHostPID := 0
	podHostIPC := 0
	podRunAsRoot := 0
	for _, pod := range data.Workloads.Pods {
		for _, c := range pod.Containers {
			if c.SecurityContext.Privileged {
				podPrivileged++
				break
			}
			if c.SecurityContext.RunAsUser != nil && *c.SecurityContext.RunAsUser == 0 {
				podRunAsRoot++
				break
			}
		}
		if pod.SecurityContext.HostNetwork {
			podHostNetwork++
		}
		if pod.SecurityContext.HostPID {
			podHostPID++
		}
		if pod.SecurityContext.HostIPC {
			podHostIPC++
		}
	}
	podLabels := []string{"Total Pods", "Privileged", "Host Network", "Host PID", "Host IPC", "RunAsRoot"}
	podCounts := []int{podTotal, podPrivileged, podHostNetwork, podHostPID, podHostIPC, podRunAsRoot}
	r.excel.SetCellValue(sheet, "E12", "Pod Security Summary")
	r.excel.MergeCell(sheet, "E12", "F12")
	r.excel.SetCellStyle(sheet, "E12", "F12", r.sectionStyle)
	r.excel.SetCellValue(sheet, "E13", "Type")
	r.excel.SetCellValue(sheet, "F13", "Count")
	r.excel.SetCellStyle(sheet, "E13", "E13", r.headerStyle)
	r.excel.SetCellStyle(sheet, "F13", "F13", r.headerStyle)
	for i, label := range podLabels {
		row := podTableStart + 2 + i
		r.excel.SetCellValue(sheet, fmt.Sprintf("E%d", row), label)
		r.excel.SetCellValue(sheet, fmt.Sprintf("F%d", row), podCounts[i])
		style := r.contentStyle
		if i%2 == 1 {
			style = r.altRowStyle
		}
		r.excel.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), style)
	}
	// Insert Pod Security chart below the table
	r.excel.AddChart(sheet, "E21", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{{
			Name:       "Pod Security",
			Categories: fmt.Sprintf("%s!$E$15:$E$20", sheet),
			Values:     fmt.Sprintf("%s!$F$15:$F$20", sheet),
		}},
		Title:  []excelize.RichTextRun{{Text: "Pod Configurations"}},
		Legend: excelize.ChartLegend{Position: "top"},
	})

	// Auto-fit columns
	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generatePods(workloads models.WorkloadAssessment) error {
	sheet := "Pods"

	// Set headers with security configurations
	headers := []string{
		"Name", "Namespace", "Node", "Service Account",
		"Privileged", "Host Network", "Host PID", "Host IPC",
		"Run As User", "Run As Non Root", "Auto Mount SA Token",
		"Container Names", "Container Images", "Capabilities",
		"Resources", "Sysctls", "Environment Variables",
		"Created At", "Labels",
	}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}

	// Add pod data
	for i, pod := range workloads.Pods {
		row := i + 2
		containerNames := make([]string, 0)
		imageNames := make([]string, 0)
		capabilities := make([]string, 0)
		resourceInfo := make([]string, 0)
		envVars := make([]string, 0)
		hasPrivileged := false
		runAsUser := int64(0)
		runAsNonRoot := false
		automountServiceAccountToken := true // default is true in K8s

		// Get pod-level security context values
		if pod.SecurityContext.RunAsUser != nil {
			runAsUser = *pod.SecurityContext.RunAsUser
		}

		// Collect container-level information
		for _, container := range pod.Containers {
			containerNames = append(containerNames, container.Name)
			imageNames = append(imageNames, container.Image)

			// Collect security context information
			if len(container.SecurityContext.Capabilities) > 0 {
				capabilities = append(capabilities, container.SecurityContext.Capabilities...)
			}

			if container.SecurityContext.Privileged {
				hasPrivileged = true
			}

			if container.SecurityContext.RunAsNonRoot != nil {
				runAsNonRoot = *container.SecurityContext.RunAsNonRoot
			}

			// Format resource information
			if container.Resources.Limits.CPU != "" {
				resourceInfo = append(resourceInfo, fmt.Sprintf("%s: CPU limit %s", container.Name, container.Resources.Limits.CPU))
			}
			if container.Resources.Limits.Memory != "" {
				resourceInfo = append(resourceInfo, fmt.Sprintf("%s: Memory limit %s", container.Name, container.Resources.Limits.Memory))
			}
			if container.Resources.Requests.CPU != "" {
				resourceInfo = append(resourceInfo, fmt.Sprintf("%s: CPU request %s", container.Name, container.Resources.Requests.CPU))
			}
			if container.Resources.Requests.Memory != "" {
				resourceInfo = append(resourceInfo, fmt.Sprintf("%s: Memory request %s", container.Name, container.Resources.Requests.Memory))
			}

			// Collect environment variables (excluding secrets)
			for _, env := range container.EnvVars {
				if !strings.Contains(strings.ToLower(env), "secret") &&
					!strings.Contains(strings.ToLower(env), "password") &&
					!strings.Contains(strings.ToLower(env), "key") {
					envVars = append(envVars, env)
				}
			}
		}

		// If pod spec has automountServiceAccountToken set, override default
		if pod.AutomountServiceAccountToken != nil {
			automountServiceAccountToken = *pod.AutomountServiceAccountToken
		}

		values := []interface{}{
			pod.Name,
			pod.Namespace,
			pod.NodeName,
			pod.ServiceAccount,
			hasPrivileged,
			pod.SecurityContext.HostNetwork,
			pod.SecurityContext.HostPID,
			pod.SecurityContext.HostIPC,
			runAsUser,
			runAsNonRoot,
			automountServiceAccountToken,
			strings.Join(containerNames, "\n"),
			strings.Join(imageNames, "\n"),
			strings.Join(capabilities, ", "),
			strings.Join(resourceInfo, "\n"),
			"N/A", // sysctls - not directly available in the model
			strings.Join(envVars, "\n"),
			pod.CreatedAt,
			r.formatLabels(pod.Labels),
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)

			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}

	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateNodes(nodes []models.NodeInfo) error {
	sheet := "Nodes"
	headers := []string{"Name", "Version", "Architecture", "OS", "Container Runtime", "CPU", "Memory", "Ready", "Labels"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	for i, node := range nodes {
		row := i + 2
		values := []interface{}{
			node.Name,
			node.Version,
			node.Architecture,
			node.OS,
			node.ContainerRuntime,
			node.CPU,
			node.Memory,
			node.Ready,
			r.formatLabels(node.Labels),
		}
		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}
	r.autoFitColumns(sheet)
	return nil
}

func (r *Report) generateNamespaces(namespaces []models.NamespaceInfo) error {
	sheet := "Namespaces"
	headers := []string{"Name", "Status", "Created At", "Labels"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		r.excel.SetCellValue(sheet, cell, header)
		r.excel.SetCellStyle(sheet, cell, cell, r.headerStyle)
	}
	for i, ns := range namespaces {
		row := i + 2
		values := []interface{}{
			ns.Name,
			ns.Status,
			ns.CreatedAt,
			r.formatLabels(ns.Labels),
		}
		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			r.excel.SetCellValue(sheet, cell, value)
			style := r.contentStyle
			if row%2 == 0 {
				style = r.altRowStyle
			}
			r.excel.SetCellStyle(sheet, cell, cell, style)
		}
	}
	r.autoFitColumns(sheet)
	return nil
}
