package dashboard

import (
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/database"
	"github.com/kgermando/dentic-support-api/models"
	"github.com/kgermando/dentic-support-api/utils"
)

// ─────────────────────────────────────────────────────────────────────────────
// Types graphiques — compatibles Chart.js / ApexCharts
// ─────────────────────────────────────────────────────────────────────────────

type PieChartData struct {
	Labels []string  `json:"labels"`
	Series []float64 `json:"series"`
	Colors []string  `json:"colors"`
}

type LineDataset struct {
	Name  string    `json:"name"`
	Data  []float64 `json:"data"`
	Color string    `json:"color"`
}

type LineChartData struct {
	Labels   []string      `json:"labels"`
	Datasets []LineDataset `json:"datasets"`
}

type BarDataset struct {
	Name  string    `json:"name"`
	Data  []float64 `json:"data"`
	Color string    `json:"color"`
}

type BarChartData struct {
	Labels   []string     `json:"labels"`
	Datasets []BarDataset `json:"datasets"`
}

type RadarDataset struct {
	Name  string    `json:"name"`
	Data  []float64 `json:"data"`
	Color string    `json:"color"`
}

type RadarChartData struct {
	Labels   []string       `json:"labels"`
	Datasets []RadarDataset `json:"datasets"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Palettes de couleurs
// ─────────────────────────────────────────────────────────────────────────────

var (
	colorsStatutTicket = []string{"#F59E0B", "#3B82F6", "#10B981", "#6B7280"}
	colorsStatutTache  = []string{"#F97316", "#3B82F6", "#10B981"}
	colorsRoles        = []string{"#8B5CF6", "#EC4899", "#14B8A6", "#F59E0B"}
	colorsTranches     = []string{"#06B6D4", "#3B82F6", "#8B5CF6", "#EC4899", "#F97316"}
	colorsPalette      = []string{
		"#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6",
		"#14B8A6", "#F97316", "#EC4899", "#6366F1", "#84CC16",
	}
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers internes
// ─────────────────────────────────────────────────────────────────────────────

type monthlyRow struct {
	Month string
	Count int64
}

type categoryRow struct {
	Category string
	Count    int64
}

// color retourne une couleur de la palette par index cyclique
func paletteColor(i int) string {
	return colorsPalette[i%len(colorsPalette)]
}

// mergedMonthLabels fusionne et trie les mois de plusieurs séries
func mergedMonthLabels(series ...[]monthlyRow) []string {
	set := make(map[string]struct{})
	for _, s := range series {
		for _, r := range s {
			set[r.Month] = struct{}{}
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// rowsToMap convertit []monthlyRow en map[month]count
func rowsToMap(rows []monthlyRow) map[string]float64 {
	m := make(map[string]float64, len(rows))
	for _, r := range rows {
		m[r.Month] = float64(r.Count)
	}
	return m
}

// ─────────────────────────────────────────────────────────────────────────────
// POINT D'ENTRÉE UNIQUE — filtre automatique sur le rôle de l'agent connecté
// GET /api/dashboard?token=<jwt>
// ─────────────────────────────────────────────────────────────────────────────

func GetDashboard(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token requis",
		})
	}

	agentUUID, err := utils.VerifyJwt(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token invalide",
		})
	}

	db := database.DB
	var agent models.Agent
	if err := db.Where("uuid = ?", agentUUID).
		Preload("Bureau").
		Preload("Direction").
		First(&agent).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Agent introuvable",
		})
	}

	if !agent.Status {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Compte désactivé",
		})
	}

	switch agent.Role {
	case "SuperAdmin", "Directeur":
		return directeurDashboard(c, agent)
	case "Chef du bureau":
		return chefBureauDashboard(c, agent)
	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Accès refusé : rôle non autorisé (" + agent.Role + ")",
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// DASHBOARD DIRECTEUR GÉNÉRAL
// ─────────────────────────────────────────────────────────────────────────────

func directeurDashboard(c *fiber.Ctx, agent models.Agent) error {
	db := database.DB

	// ── KPIs ─────────────────────────────────────────────────────────────────
	var totalDirections, totalBureaux, totalAgents, totalAgentsActifs int64
	var totalDemandeurs, totalTickets, totalTaches, totalTeams int64

	db.Model(&models.Direction{}).Count(&totalDirections)
	db.Model(&models.Bureau{}).Count(&totalBureaux)
	db.Model(&models.Agent{}).Count(&totalAgents)
	db.Model(&models.Agent{}).Where("status = true").Count(&totalAgentsActifs)
	db.Model(&models.Demandeur{}).Count(&totalDemandeurs)
	db.Model(&models.Ticket{}).Count(&totalTickets)
	db.Model(&models.Tache{}).Count(&totalTaches)
	db.Model(&models.Team{}).Count(&totalTeams)

	var resolus, ticketsEnRetard int64
	db.Model(&models.Ticket{}).Where("statut IN ?", []string{"Résolu", "Fermé"}).Count(&resolus)
	db.Model(&models.Ticket{}).Where("statut = 'Ouvert' AND created_at < ?", time.Now().AddDate(0, 0, -3)).Count(&ticketsEnRetard)

	tauxResolution := 0.0
	if totalTickets > 0 {
		tauxResolution = float64(resolus) / float64(totalTickets) * 100
	}

	kpis := fiber.Map{
		"total_directions":    totalDirections,
		"total_bureaux":       totalBureaux,
		"total_agents":        totalAgents,
		"total_agents_actifs": totalAgentsActifs,
		"total_demandeurs":    totalDemandeurs,
		"total_tickets":       totalTickets,
		"total_taches":        totalTaches,
		"total_teams":         totalTeams,
		"taux_resolution":     tauxResolution,
		"tickets_en_retard":   ticketsEnRetard,
	}

	// ── PieChart 1 : Tickets par statut ──────────────────────────────────────
	ticketsStatutPie := PieChartData{
		Labels: []string{"Ouvert", "En cours", "Résolu", "Fermé"},
		Colors: colorsStatutTicket,
	}
	for _, s := range ticketsStatutPie.Labels {
		var cnt int64
		db.Model(&models.Ticket{}).Where("statut = ?", s).Count(&cnt)
		ticketsStatutPie.Series = append(ticketsStatutPie.Series, float64(cnt))
	}

	// ── PieChart 2 : Tâches par statut ───────────────────────────────────────
	tachesStatutPie := PieChartData{
		Labels: []string{"En attente", "En cours", "Terminé"},
		Colors: colorsStatutTache,
	}
	for _, s := range tachesStatutPie.Labels {
		var cnt int64
		db.Model(&models.Tache{}).Where("statut = ?", s).Count(&cnt)
		tachesStatutPie.Series = append(tachesStatutPie.Series, float64(cnt))
	}

	// ── PieChart 3 : Agents par rôle ─────────────────────────────────────────
	agentsRolePie := PieChartData{
		Labels: []string{"Directeur", "Secretaire", "Chef du bureau", "Agent"},
		Colors: colorsRoles,
	}
	for _, r := range agentsRolePie.Labels {
		var cnt int64
		db.Model(&models.Agent{}).Where("role = ?", r).Count(&cnt)
		agentsRolePie.Series = append(agentsRolePie.Series, float64(cnt))
	}

	// ── PieChart 4 : Agents par tranche d'âge ────────────────────────────────
	type trancheRow struct {
		TranchAge string `gorm:"column:tranch_age"`
		Count     int64
	}
	var trRows []trancheRow
	db.Model(&models.Agent{}).
		Select("tranch_age, COUNT(*) as count").
		Group("tranch_age").Order("tranch_age ASC").Scan(&trRows)

	agentsTranchePie := PieChartData{}
	for i, r := range trRows {
		agentsTranchePie.Labels = append(agentsTranchePie.Labels, r.TranchAge)
		agentsTranchePie.Series = append(agentsTranchePie.Series, float64(r.Count))
		c := "#94A3B8"
		if i < len(colorsTranches) {
			c = colorsTranches[i]
		}
		agentsTranchePie.Colors = append(agentsTranchePie.Colors, c)
	}

	// ── BarChart : Tickets par catégorie ─────────────────────────────────────
	var catRows []categoryRow
	db.Model(&models.Ticket{}).
		Select("category, COUNT(*) as count").
		Group("category").Order("count DESC").Scan(&catRows)

	catBar := BarChartData{}
	catValues := []float64{}
	for _, r := range catRows {
		catBar.Labels = append(catBar.Labels, r.Category)
		catValues = append(catValues, float64(r.Count))
	}
	catBar.Datasets = []BarDataset{{Name: "Tickets", Data: catValues, Color: "#3B82F6"}}

	// ── StackedBarChart : Tickets par bureau, décomposé par statut ───────────
	var bureaux []models.Bureau
	db.Order("name ASC").Find(&bureaux)

	bureauStacked := BarChartData{}
	sOuverts, sEnCours, sResolus, sFermes := []float64{}, []float64{}, []float64{}, []float64{}
	for _, b := range bureaux {
		bureauStacked.Labels = append(bureauStacked.Labels, b.Name)
		var o, e, r, f int64
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut=?", b.UUID, "Ouvert").Count(&o)
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut=?", b.UUID, "En cours").Count(&e)
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut=?", b.UUID, "Résolu").Count(&r)
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut=?", b.UUID, "Fermé").Count(&f)
		sOuverts = append(sOuverts, float64(o))
		sEnCours = append(sEnCours, float64(e))
		sResolus = append(sResolus, float64(r))
		sFermes = append(sFermes, float64(f))
	}
	bureauStacked.Datasets = []BarDataset{
		{Name: "Ouvert", Data: sOuverts, Color: colorsStatutTicket[0]},
		{Name: "En cours", Data: sEnCours, Color: colorsStatutTicket[1]},
		{Name: "Résolu", Data: sResolus, Color: colorsStatutTicket[2]},
		{Name: "Fermé", Data: sFermes, Color: colorsStatutTicket[3]},
	}

	// ── BarChart groupé : Directions comparaison ──────────────────────────────
	var directions []models.Direction
	db.Order("name ASC").Find(&directions)

	dirBar := BarChartData{}
	dAgents, dBureaux, dTickets := []float64{}, []float64{}, []float64{}
	for _, d := range directions {
		dirBar.Labels = append(dirBar.Labels, d.Name)
		var a, b, t int64
		db.Model(&models.Agent{}).Where("direction_uuid=?", d.UUID).Count(&a)
		db.Model(&models.Bureau{}).Where("direction_uuid=?", d.UUID).Count(&b)
		var bUUIDs []string
		db.Model(&models.Bureau{}).Where("direction_uuid=?", d.UUID).Pluck("uuid", &bUUIDs)
		if len(bUUIDs) > 0 {
			db.Model(&models.Ticket{}).Where("bureau_uuid IN ?", bUUIDs).Count(&t)
		}
		dAgents = append(dAgents, float64(a))
		dBureaux = append(dBureaux, float64(b))
		dTickets = append(dTickets, float64(t))
	}
	dirBar.Datasets = []BarDataset{
		{Name: "Agents", Data: dAgents, Color: "#8B5CF6"},
		{Name: "Bureaux", Data: dBureaux, Color: "#14B8A6"},
		{Name: "Tickets", Data: dTickets, Color: "#F59E0B"},
	}

	// ── LineChart : Évolution mensuelle tickets + tâches (12 mois) ───────────
	var tkMonths, taMonths []monthlyRow
	db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',created_at),'YYYY-MM') AS month, COUNT(*) AS count
		FROM tickets WHERE deleted_at IS NULL AND created_at>=NOW()-INTERVAL '12 months'
		GROUP BY month ORDER BY month ASC`).Scan(&tkMonths)
	db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',created_at),'YYYY-MM') AS month, COUNT(*) AS count
		FROM taches WHERE deleted_at IS NULL AND created_at>=NOW()-INTERVAL '12 months'
		GROUP BY month ORDER BY month ASC`).Scan(&taMonths)

	allMonths := mergedMonthLabels(tkMonths, taMonths)
	tkMap, taMap := rowsToMap(tkMonths), rowsToMap(taMonths)
	evolutionLine := LineChartData{Labels: allMonths}
	tkData, taData := []float64{}, []float64{}
	for _, m := range allMonths {
		tkData = append(tkData, tkMap[m])
		taData = append(taData, taMap[m])
	}
	evolutionLine.Datasets = []LineDataset{
		{Name: "Tickets", Data: tkData, Color: "#3B82F6"},
		{Name: "Tâches", Data: taData, Color: "#10B981"},
	}

	// ── LineChart : Tickets créés vs résolus par mois ─────────────────────────
	var resolvedMonths []monthlyRow
	db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',updated_at),'YYYY-MM') AS month, COUNT(*) AS count
		FROM tickets WHERE deleted_at IS NULL AND statut IN ('Résolu','Fermé')
		AND updated_at>=NOW()-INTERVAL '12 months' GROUP BY month ORDER BY month ASC`).Scan(&resolvedMonths)

	resMap := rowsToMap(resolvedMonths)
	createdVsResolved := LineChartData{}
	cData, rData := []float64{}, []float64{}
	for _, r := range tkMonths {
		createdVsResolved.Labels = append(createdVsResolved.Labels, r.Month)
		cData = append(cData, float64(r.Count))
		rData = append(rData, resMap[r.Month])
	}
	createdVsResolved.Datasets = []LineDataset{
		{Name: "Créés", Data: cData, Color: "#F59E0B"},
		{Name: "Résolus", Data: rData, Color: "#10B981"},
	}

	// ── RadarChart : Taux de résolution par bureau ────────────────────────────
	radarBureau := RadarChartData{}
	radarScores := []float64{}
	for _, b := range bureaux {
		radarBureau.Labels = append(radarBureau.Labels, b.Name)
		var tot, res int64
		db.Model(&models.Ticket{}).Where("bureau_uuid=?", b.UUID).Count(&tot)
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut IN ?", b.UUID, []string{"Résolu", "Fermé"}).Count(&res)
		score := 0.0
		if tot > 0 {
			score = float64(res) / float64(tot) * 100
		}
		radarScores = append(radarScores, score)
	}
	radarBureau.Datasets = []RadarDataset{
		{Name: "Taux de résolution (%)", Data: radarScores, Color: "#6366F1"},
	}

	// ── Tableaux récents ──────────────────────────────────────────────────────
	var recentTickets []models.Ticket
	db.Preload("Demandeur").Preload("Bureau").Order("created_at DESC").Limit(10).Find(&recentTickets)

	var recentAgents []models.Agent
	db.Preload("Direction").Preload("Bureau").Order("created_at DESC").Limit(10).Find(&recentAgents)

	return c.JSON(fiber.Map{
		"status":  "success",
		"role":    agent.Role,
		"message": "Dashboard Directeur Général",
		"data": fiber.Map{
			"agent":          agent,
			"kpis":           kpis,
			"recent_tickets": recentTickets,
			"recent_agents":  recentAgents,
			"charts": fiber.Map{
				// Pie / Donut
				"pie_tickets_par_statut":     ticketsStatutPie,
				"pie_taches_par_statut":      tachesStatutPie,
				"pie_agents_par_role":        agentsRolePie,
				"pie_agents_par_tranche_age": agentsTranchePie,
				// Bar
				"bar_tickets_par_categorie":      catBar,
				"bar_direction_comparaison":      dirBar,
				"stacked_bar_tickets_par_bureau": bureauStacked,
				// Line
				"line_evolution_tickets_taches": evolutionLine,
				"line_crees_vs_resolus":         createdVsResolved,
				// Radar
				"radar_taux_resolution_bureaux": radarBureau,
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// DASHBOARD CHEF DE BUREAU
// ─────────────────────────────────────────────────────────────────────────────

func chefBureauDashboard(c *fiber.Ctx, agent models.Agent) error {
	bureauUUID := agent.BureauUUID
	db := database.DB

	var bureau models.Bureau
	db.Where("uuid=?", bureauUUID).Preload("Direction").First(&bureau)

	// Agents du bureau
	var agentUUIDs []string
	db.Model(&models.Agent{}).Where("bureau_uuid=?", bureauUUID).Pluck("uuid", &agentUUIDs)

	// ── KPIs ─────────────────────────────────────────────────────────────────
	var totalAgents, totalAgentsActifs, totalTickets, totalTaches, totalMembres int64
	db.Model(&models.Agent{}).Where("bureau_uuid=?", bureauUUID).Count(&totalAgents)
	db.Model(&models.Agent{}).Where("bureau_uuid=? AND status=true", bureauUUID).Count(&totalAgentsActifs)
	db.Model(&models.Ticket{}).Where("bureau_uuid=?", bureauUUID).Count(&totalTickets)
	if len(agentUUIDs) > 0 {
		db.Model(&models.Tache{}).Where("agent_uuid IN ?", agentUUIDs).Count(&totalTaches)
	}
	db.Model(&models.TeamJoin{}).Where("bureau_uuid=?", bureauUUID).Count(&totalMembres)

	var resolus, ticketsEnRetard, ticketsSansTache int64
	db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut IN ?", bureauUUID, []string{"Résolu", "Fermé"}).Count(&resolus)
	db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut='Ouvert' AND created_at<?", bureauUUID, time.Now().AddDate(0, 0, -3)).Count(&ticketsEnRetard)

	var ticketsAvecTache []string
	db.Model(&models.Tache{}).Select("DISTINCT ticket_uuid").Scan(&ticketsAvecTache)
	if len(ticketsAvecTache) > 0 {
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND uuid NOT IN ? AND statut NOT IN ?", bureauUUID, ticketsAvecTache, []string{"Résolu", "Fermé"}).Count(&ticketsSansTache)
	} else {
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut NOT IN ?", bureauUUID, []string{"Résolu", "Fermé"}).Count(&ticketsSansTache)
	}

	tauxResolution := 0.0
	if totalTickets > 0 {
		tauxResolution = float64(resolus) / float64(totalTickets) * 100
	}

	kpis := fiber.Map{
		"bureau_uuid":           bureauUUID,
		"bureau_name":           bureau.Name,
		"direction_name":        bureau.Direction.Name,
		"total_agents":          totalAgents,
		"total_agents_actifs":   totalAgentsActifs,
		"total_tickets":         totalTickets,
		"total_taches":          totalTaches,
		"total_membres_equipes": totalMembres,
		"taux_resolution":       tauxResolution,
		"tickets_en_retard":     ticketsEnRetard,
		"tickets_sans_tache":    ticketsSansTache,
	}

	// ── PieChart 1 : Tickets par statut ──────────────────────────────────────
	ticketsStatutPie := PieChartData{
		Labels: []string{"Ouvert", "En cours", "Résolu", "Fermé"},
		Colors: colorsStatutTicket,
	}
	for _, s := range ticketsStatutPie.Labels {
		var cnt int64
		db.Model(&models.Ticket{}).Where("bureau_uuid=? AND statut=?", bureauUUID, s).Count(&cnt)
		ticketsStatutPie.Series = append(ticketsStatutPie.Series, float64(cnt))
	}

	// ── PieChart 2 : Tickets par catégorie ───────────────────────────────────
	var catRows []categoryRow
	db.Model(&models.Ticket{}).
		Select("category, COUNT(*) as count").
		Where("bureau_uuid=?", bureauUUID).
		Group("category").Order("count DESC").Scan(&catRows)

	ticketsCatPie := PieChartData{}
	for i, r := range catRows {
		ticketsCatPie.Labels = append(ticketsCatPie.Labels, r.Category)
		ticketsCatPie.Series = append(ticketsCatPie.Series, float64(r.Count))
		ticketsCatPie.Colors = append(ticketsCatPie.Colors, paletteColor(i))
	}

	// ── PieChart 3 : Tâches par statut ───────────────────────────────────────
	tachesStatutPie := PieChartData{
		Labels: []string{"En attente", "En cours", "Terminé"},
		Colors: colorsStatutTache,
		Series: []float64{0, 0, 0},
	}
	if len(agentUUIDs) > 0 {
		tachesStatutPie.Series = []float64{}
		for _, s := range tachesStatutPie.Labels {
			var cnt int64
			db.Model(&models.Tache{}).Where("agent_uuid IN ? AND statut=?", agentUUIDs, s).Count(&cnt)
			tachesStatutPie.Series = append(tachesStatutPie.Series, float64(cnt))
		}
	}

	// ── StackedBarChart : Performance agents (tâches) ─────────────────────────
	var agentsList []models.Agent
	db.Where("bureau_uuid=?", bureauUUID).Order("fullname ASC").Find(&agentsList)

	agentPerfBar := BarChartData{}
	pTermine, pEnCours, pAttente := []float64{}, []float64{}, []float64{}
	for _, a := range agentsList {
		agentPerfBar.Labels = append(agentPerfBar.Labels, a.Fullname)
		var t, e, at int64
		db.Model(&models.Tache{}).Where("agent_uuid=? AND statut='Terminé'", a.UUID).Count(&t)
		db.Model(&models.Tache{}).Where("agent_uuid=? AND statut='En cours'", a.UUID).Count(&e)
		db.Model(&models.Tache{}).Where("agent_uuid=? AND statut='En attente'", a.UUID).Count(&at)
		pTermine = append(pTermine, float64(t))
		pEnCours = append(pEnCours, float64(e))
		pAttente = append(pAttente, float64(at))
	}
	agentPerfBar.Datasets = []BarDataset{
		{Name: "Terminé", Data: pTermine, Color: colorsStatutTache[2]},
		{Name: "En cours", Data: pEnCours, Color: colorsStatutTache[1]},
		{Name: "En attente", Data: pAttente, Color: colorsStatutTache[0]},
	}

	// ── BarChart : Tickets traités par agent ──────────────────────────────────
	agentTicketsBar := BarChartData{}
	agentTicketsCounts := []float64{}
	for _, a := range agentsList {
		agentTicketsBar.Labels = append(agentTicketsBar.Labels, a.Fullname)
		var cnt int64
		db.Model(&models.Tache{}).Select("COUNT(DISTINCT ticket_uuid)").Where("agent_uuid=?", a.UUID).Scan(&cnt)
		agentTicketsCounts = append(agentTicketsCounts, float64(cnt))
	}
	agentTicketsBar.Datasets = []BarDataset{
		{Name: "Tickets traités", Data: agentTicketsCounts, Color: "#6366F1"},
	}

	// ── LineChart : Évolution mensuelle tickets + tâches (bureau) ────────────
	var tkMonths []monthlyRow
	db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',created_at),'YYYY-MM') AS month, COUNT(*) AS count
		FROM tickets WHERE deleted_at IS NULL AND bureau_uuid=?
		AND created_at>=NOW()-INTERVAL '12 months' GROUP BY month ORDER BY month ASC`, bureauUUID).Scan(&tkMonths)

	var taMonths []monthlyRow
	if len(agentUUIDs) > 0 {
		db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',created_at),'YYYY-MM') AS month, COUNT(*) AS count
			FROM taches WHERE deleted_at IS NULL AND agent_uuid=ANY(?)
			AND created_at>=NOW()-INTERVAL '12 months' GROUP BY month ORDER BY month ASC`, agentUUIDs).Scan(&taMonths)
	}

	allMonths := mergedMonthLabels(tkMonths, taMonths)
	tkMap, taMap := rowsToMap(tkMonths), rowsToMap(taMonths)
	evolutionLine := LineChartData{Labels: allMonths}
	tkData, taData := []float64{}, []float64{}
	for _, m := range allMonths {
		tkData = append(tkData, tkMap[m])
		taData = append(taData, taMap[m])
	}
	evolutionLine.Datasets = []LineDataset{
		{Name: "Tickets", Data: tkData, Color: "#3B82F6"},
		{Name: "Tâches", Data: taData, Color: "#10B981"},
	}

	// ── LineChart : Créés vs résolus (bureau) ─────────────────────────────────
	var resolvedMonths []monthlyRow
	db.Raw(`SELECT TO_CHAR(DATE_TRUNC('month',updated_at),'YYYY-MM') AS month, COUNT(*) AS count
		FROM tickets WHERE deleted_at IS NULL AND bureau_uuid=?
		AND statut IN ('Résolu','Fermé') AND updated_at>=NOW()-INTERVAL '12 months'
		GROUP BY month ORDER BY month ASC`, bureauUUID).Scan(&resolvedMonths)

	resMap := rowsToMap(resolvedMonths)
	createdVsResolved := LineChartData{}
	cData, rData := []float64{}, []float64{}
	for _, r := range tkMonths {
		createdVsResolved.Labels = append(createdVsResolved.Labels, r.Month)
		cData = append(cData, float64(r.Count))
		rData = append(rData, resMap[r.Month])
	}
	createdVsResolved.Datasets = []LineDataset{
		{Name: "Créés", Data: cData, Color: "#F59E0B"},
		{Name: "Résolus", Data: rData, Color: "#10B981"},
	}

	// ── RadarChart : Profil des agents ────────────────────────────────────────
	radarAgents := RadarChartData{
		Labels: []string{"Tâches terminées", "Tâches en cours", "Tickets traités"},
	}
	for i, a := range agentsList {
		var termine, enCours, traites int64
		db.Model(&models.Tache{}).Where("agent_uuid=? AND statut='Terminé'", a.UUID).Count(&termine)
		db.Model(&models.Tache{}).Where("agent_uuid=? AND statut='En cours'", a.UUID).Count(&enCours)
		db.Model(&models.Tache{}).Select("COUNT(DISTINCT ticket_uuid)").Where("agent_uuid=?", a.UUID).Scan(&traites)
		radarAgents.Datasets = append(radarAgents.Datasets, RadarDataset{
			Name:  a.Fullname,
			Data:  []float64{float64(termine), float64(enCours), float64(traites)},
			Color: paletteColor(i),
		})
	}

	// ── Tableau récent ────────────────────────────────────────────────────────
	var recentTickets []models.Ticket
	db.Where("bureau_uuid=?", bureauUUID).
		Preload("Demandeur").Order("created_at DESC").Limit(10).Find(&recentTickets)

	// ── Équipes ───────────────────────────────────────────────────────────────
	var teamJoins []models.TeamJoin
	db.Where("bureau_uuid=?", bureauUUID).Preload("Team").Preload("Agent").Find(&teamJoins)

	type teamEntry struct {
		TeamUUID string      `json:"team_uuid"`
		TeamName string      `json:"team_name"`
		Membres  []fiber.Map `json:"membres"`
	}
	teamsMap := make(map[string]*teamEntry)
	for _, tj := range teamJoins {
		if _, ok := teamsMap[tj.TeamUUID]; !ok {
			teamsMap[tj.TeamUUID] = &teamEntry{TeamUUID: tj.TeamUUID, TeamName: tj.Team.Name}
		}
		teamsMap[tj.TeamUUID].Membres = append(teamsMap[tj.TeamUUID].Membres, fiber.Map{
			"agent_uuid": tj.AgentUUID,
			"fullname":   tj.Agent.Fullname,
			"role":       tj.Agent.Role,
		})
	}
	var teamsBreakdown []teamEntry
	for _, v := range teamsMap {
		teamsBreakdown = append(teamsBreakdown, *v)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"role":    agent.Role,
		"message": "Dashboard Chef de Bureau — " + bureau.Name,
		"data": fiber.Map{
			"agent":           agent,
			"kpis":            kpis,
			"recent_tickets":  recentTickets,
			"teams_breakdown": teamsBreakdown,
			"charts": fiber.Map{
				// Pie / Donut
				"pie_tickets_par_statut":    ticketsStatutPie,
				"pie_tickets_par_categorie": ticketsCatPie,
				"pie_taches_par_statut":     tachesStatutPie,
				// Bar
				"stacked_bar_performance_agents": agentPerfBar,
				"bar_tickets_traites_par_agent":  agentTicketsBar,
				// Line
				"line_evolution_tickets_taches": evolutionLine,
				"line_crees_vs_resolus":         createdVsResolved,
				// Radar
				"radar_profil_agents": radarAgents,
			},
		},
	})
}
