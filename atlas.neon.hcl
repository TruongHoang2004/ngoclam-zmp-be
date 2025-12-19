data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-gorm", "load",
    "--path", "./internal/infrastructure/persistence/model",
    "--dialect", "postgres"
  ]
}

env "gorm" {
  # Nguồn schema lấy từ GORM
  src = data.external_schema.gorm.url

  # Database dev (chỉ dùng để Atlas diff, không phải DB production)
  dev = "postgres://postgres:truonghoang2004@localhost:5432/empty?sslmode=disable"

  # Real DB (Neon)
  url = "postgresql://neondb_owner:npg_KwLRlgn3pd5D@ep-plain-resonance-a1adxptm-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

  migration {
    dir = "file://internal/infrastructure/database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
