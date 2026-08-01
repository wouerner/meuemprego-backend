package database

import (
	"log"

	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/gorm"
)

// SeedData popula o banco com dados de exemplo (idempotente: só insere se a tabela estiver vazia).
func SeedData(db *gorm.DB) error {
	if err := seedCandidates(db); err != nil {
		return err
	}
	if err := seedHunters(db); err != nil {
		return err
	}
	return nil
}

func seedCandidates(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.Candidate{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	candidates := []domain.Candidate{
		{
			Name:                 "Carlos Eduardo",
			CPF:                  "529.982.247-25",
			Email:                "carlos.eduardo@email.com",
			Avatar:               "https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=250&q=80",
			Headline:             "Engenheiro de Software Senior | Fullstack (Vue.js, Node.js)",
			Seniority:            "Senior",
			Area:                 "Tecnologia da Informação",
			CareerGoal:           "Transição para Tech Lead ou Engenheiro Staff em empresa global de tecnologia.",
			ProfessionalMoment:   "Aberto a Propostas",
			RequestHunterContact: true,
			LGPDConsent:          true,
			IsApproved:           true,
			LinkedInURL:          "https://www.linkedin.com/in/carlos-eduardo-demo",
			WhatsAppNumber:       "5511998765432",
		},
		{
			Name:                 "Mariana Silva",
			CPF:                  "111.444.777-35",
			Email:                "mariana.silva@email.com",
			Avatar:               "https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=250&q=80",
			Headline:             "Product Manager Pleno | Experiência em Fintechs & B2B SaaS",
			Seniority:            "Pleno",
			Area:                 "Produtos & Design",
			CareerGoal:           "Aceleração para cadeira de Senior Product Manager com liderança de squad.",
			ProfessionalMoment:   "Em Transição",
			RequestHunterContact: true,
			LGPDConsent:          true,
			IsApproved:           true,
			LinkedInURL:          "https://www.linkedin.com/in/mariana-silva-pm",
			WhatsAppNumber:       "5511987654321",
		},
		{
			Name:                 "Gabriel Santos",
			CPF:                  "222.333.444-55",
			Email:                "gabriel.santos@email.com",
			Avatar:               "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=250&q=80",
			Headline:             "Diretor de Vendas & Expansion | Growth B2B",
			Seniority:            "Liderança / C-Level",
			Area:                 "Vendas / Comercial",
			CareerGoal:           "Colocação executiva como VP of Sales ou Chief Revenue Officer (CRO).",
			ProfessionalMoment:   "Buscando recolocação",
			RequestHunterContact: true,
			LGPDConsent:          true,
			IsApproved:           true,
			LinkedInURL:          "https://www.linkedin.com/in/gabriel-santos-cro",
			WhatsAppNumber:       "5511976543219",
		},
		{
			Name:                 "Beatriz Lima",
			CPF:                  "333.444.555-66",
			Email:                "beatriz.lima@email.com",
			Avatar:               "https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&w=250&q=80",
			Headline:             "Cientista de Dados Pleno | Python, Machine Learning & SQL",
			Seniority:            "Pleno",
			Area:                 "Data & Analytics",
			CareerGoal:           "Conquistar primeira oportunidade internacional remota em IA e Data Science.",
			ProfessionalMoment:   "Ativo",
			RequestHunterContact: true,
			LGPDConsent:          true,
			IsApproved:           true,
			LinkedInURL:          "https://www.linkedin.com/in/beatriz-lima-ds",
			WhatsAppNumber:       "5521998877665",
		},
		{
			Name:                 "Felipe Oliveira",
			CPF:                  "444.555.666-77",
			Email:                "felipe.oliveira@email.com",
			Avatar:               "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?auto=format&fit=crop&w=250&q=80",
			Headline:             "Especialista de RH & Business Partner (HRBP)",
			Seniority:            "Especialista",
			Area:                 "Recursos Humanos",
			CareerGoal:           "Mentoria para reformulação estratégica de currículo e entrevistas em multinacionais.",
			ProfessionalMoment:   "Aberto a Propostas",
			RequestHunterContact: true,
			LGPDConsent:          true,
			IsApproved:           true,
			LinkedInURL:          "https://www.linkedin.com/in/felipe-oliveira-hrbp",
			WhatsAppNumber:       "5531988776655",
		},
	}

	if err := db.Create(&candidates).Error; err != nil {
		return err
	}
	log.Printf("Seed: %d candidatos de exemplo inseridos", len(candidates))
	return nil
}

func seedHunters(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.Hunter{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hunters := []domain.Hunter{
		{
			Name:               "Juliana Mendes",
			CPF:                "987.654.321-09",
			Email:              "juliana.mendes@career.com",
			Avatar:             "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&w=250&q=80",
			Headline:           "Executive Headhunter & Coach de Carreira Tech",
			Bio:                "Especialista em recolocação de executivos de TI, Tech Leads e Product Managers em grandes tech companies e startups globais.",
			Specialties:        `["Tecnologia da Informação","Produto","Engenharia de Software"]`,
			SenioritiesServed:  `["Senior","Especialista","Liderança / C-Level"]`,
			ServiceModel:       "Assessoria Completa",
			Status:             domain.HunterStatusAprovado,
			Rating:             4.9,
			TotalContactsCount: 142,
			LinkedInURL:        "https://www.linkedin.com/in/juliana-mendes-headhunter",
			WhatsAppNumber:     "5511988887777",
		},
		{
			Name:               "Roberto Andrade",
			CPF:                "222.333.444-55",
			Email:              "roberto.andrade@huntcareers.com",
			Avatar:             "https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&w=250&q=80",
			Headline:           "Job Hunter Sênior para Mercado Financeiro & Vendas",
			Bio:                "Foco total em transição de carreira para Gerentes de Contas, Diretores Comerciais e Analistas do setor bancário/fintechs.",
			Specialties:        `["Finanças","Vendas / Comercial","Gestão de Negócios"]`,
			SenioritiesServed:  `["Pleno","Senior","Liderança / C-Level"]`,
			ServiceModel:       "Mentoria de Carreira",
			Status:             domain.HunterStatusAprovado,
			Rating:             4.8,
			TotalContactsCount: 98,
			LinkedInURL:        "https://www.linkedin.com/in/roberto-andrade-hunter",
			WhatsAppNumber:     "5511977776666",
		},
		{
			Name:               "Camila Vasconcelos",
			CPF:                "333.444.555-66",
			Email:              "camila.v@jobhunter.io",
			Avatar:             "https://images.unsplash.com/photo-1580489944761-15a19d654956?auto=format&fit=crop&w=250&q=80",
			Headline:           "Especialista em Otimização de LinkedIn & Entrevistas em Inglês",
			Bio:                "Ajudo profissionais brasileiros a conquistarem vagas remotas pagas em dólares e euros através de currículos e perfis campeões.",
			Specialties:        `["Carreira Internacional","Tecnologia da Informação","Marketing"]`,
			SenioritiesServed:  `["Junior","Pleno","Senior","Especialista"]`,
			ServiceModel:       "Revisão de LinkedIn/CV",
			Status:             domain.HunterStatusAprovado,
			Rating:             5.0,
			TotalContactsCount: 215,
			LinkedInURL:        "https://www.linkedin.com/in/camila-vasconcelos-coach",
			WhatsAppNumber:     "5521999887766",
		},
		{
			Name:               "Fernando Garcia",
			CPF:                "444.555.666-77",
			Email:              "fernando.garcia@talentconsulting.com",
			Avatar:             "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?auto=format&fit=crop&w=250&q=80",
			Headline:           "Headhunter para Recursos Humanos & Operations",
			Bio:                "Conecto profissionais de RH, People Ops e Operações com empresas que valorizam cultura e diversidade.",
			Specialties:        `["Recursos Humanos","Operações & Logística"]`,
			SenioritiesServed:  `["Pleno","Senior"]`,
			ServiceModel:       "Sessão Individual",
			Status:             domain.HunterStatusPendente,
			Rating:             4.7,
			TotalContactsCount: 34,
			LinkedInURL:        "https://www.linkedin.com/in/fernando-garcia-hunter",
			WhatsAppNumber:     "5531987654321",
		},
		{
			Name:               "Luciana Rocha",
			CPF:                "555.666.777-88",
			Email:              "luciana.rocha@techhunters.com",
			Avatar:             "https://images.unsplash.com/photo-1573497019940-1c28c88b4f3e?auto=format&fit=crop&w=250&q=80",
			Headline:           "Career Strategist para Data Science & AI",
			Bio:                "Especializada no ecossistema de Ciência de Dados, Engenharia de Dados e Inteligência Artificial.",
			Specialties:        `["Data & Analytics","Tecnologia da Informação"]`,
			SenioritiesServed:  `["Senior","Especialista"]`,
			ServiceModel:       "Assessoria Completa",
			Status:             domain.HunterStatusPendente,
			Rating:             4.9,
			TotalContactsCount: 12,
			LinkedInURL:        "https://www.linkedin.com/in/luciana-rocha-ai",
			WhatsAppNumber:     "5511976543210",
		},
	}

	if err := db.Create(&hunters).Error; err != nil {
		return err
	}
	log.Printf("Seed: %d job hunters de exemplo inseridos", len(hunters))
	return nil
}
