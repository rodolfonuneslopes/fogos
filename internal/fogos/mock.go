package fogos

// MockClient returns realistic fake incidents for UI development.
// Enabled via FOGOS_MOCK=true; never used in production.
type MockClient struct{}

func (MockClient) ActiveIncidents(concelho string) ([]Incident, error) {
	if concelho == "Évora" {
		return []Incident{}, nil
	}
	return []Incident{
		{
			ID:             "mock-001",
			Concelho:       concelho,
			Freguesia:      "São Miguel",
			District:       "Castelo Branco",
			Status:         "Em Curso",
			StatusCode:     5,
			Natureza:       "Incêndio Rural",
			Date:           "01-07-2026",
			Hour:           "14:32",
			Man:            47,
			Terrain:        12,
			Aerial:         2,
			MeiosAquaticos: 0,
			DetailLocation: "São Miguel. Castelo Branco. Portugal (EN18 Km 12)",
			Extra:          "16:00 - EN18 cortada ao trânsito",
		},
		{
			ID:             "mock-002",
			Concelho:       concelho,
			Freguesia:      "Lardosa",
			District:       "Castelo Branco",
			Status:         "Em Resolução",
			StatusCode:     7,
			Natureza:       "Incêndio Rural",
			Date:           "01-07-2026",
			Hour:           "13:10",
			Man:            18,
			Terrain:        5,
			Aerial:         0,
			MeiosAquaticos: 0,
			DetailLocation: "Lardosa. Castelo Branco. Portugal",
			Extra:          "",
		},
	}, nil
}
