package database

import (
	"fmt"
	"log"
)

func (db *DB) Migrate() error {
	log.Println("Running database migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_users_table",
			sql: `
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					email TEXT UNIQUE NOT NULL,
					password_hash TEXT NOT NULL,
					needs_password_update INTEGER DEFAULT 0,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_languages_table",
			sql: `
				CREATE TABLE IF NOT EXISTS languages (
					code TEXT PRIMARY KEY,
					name TEXT NOT NULL,
					flag_url TEXT,
					active INTEGER DEFAULT 1,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_events_table",
			sql: `
				CREATE TABLE IF NOT EXISTS events (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					uuid TEXT UNIQUE NOT NULL,
					state TEXT NOT NULL DEFAULT 'draft',
					on_home INTEGER DEFAULT 0,
					title TEXT NOT NULL,
					geolink TEXT,
					type TEXT NOT NULL,
					location TEXT,
					start_date DATETIME NOT NULL,
					end_date DATETIME NOT NULL,
					image_url TEXT,
					ticket_url TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_event_translations_table",
			sql: `
				CREATE TABLE IF NOT EXISTS event_translations (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					event_id INTEGER NOT NULL,
					lang_code TEXT NOT NULL,
					description TEXT,
					FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
					FOREIGN KEY (lang_code) REFERENCES languages(code) ON DELETE CASCADE,
					UNIQUE(event_id, lang_code)
				);
			`,
		},
		{
			name: "create_geocaches_table",
			sql: `
				CREATE TABLE IF NOT EXISTS geocaches (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					gc_code TEXT UNIQUE NOT NULL,
					name TEXT NOT NULL,
					latitude REAL,
					longitude REAL,
					difficulty REAL,
					terrain REAL,
					size TEXT,
					status TEXT DEFAULT 'active',
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_messages_table",
			sql: `
				CREATE TABLE IF NOT EXISTS messages (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					state TEXT NOT NULL DEFAULT 'draft',
					priority INTEGER DEFAULT 0,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_message_translations_table",
			sql: `
				CREATE TABLE IF NOT EXISTS message_translations (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					message_id INTEGER NOT NULL,
					lang_code TEXT NOT NULL,
					title TEXT,
					content TEXT,
					FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
					FOREIGN KEY (lang_code) REFERENCES languages(code) ON DELETE CASCADE,
					UNIQUE(message_id, lang_code)
				);
			`,
		},
		{
			name: "create_static_content_table",
			sql: `
				CREATE TABLE IF NOT EXISTS static_content (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					property TEXT NOT NULL,
					lang_code TEXT NOT NULL,
					content TEXT,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (lang_code) REFERENCES languages(code) ON DELETE CASCADE,
					UNIQUE(property, lang_code)
				);
			`,
		},
		{
			name: "create_socials_table",
			sql: `
				CREATE TABLE IF NOT EXISTS socials (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					platform TEXT NOT NULL,
					url TEXT NOT NULL,
					icon TEXT,
					active INTEGER DEFAULT 1,
					sort_order INTEGER DEFAULT 0,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_contact_submissions_table",
			sql: `
				CREATE TABLE IF NOT EXISTS contact_submissions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					email TEXT NOT NULL,
					subject TEXT NOT NULL,
					message TEXT NOT NULL,
					status TEXT DEFAULT 'new',
					assigned_to INTEGER,
					last_reminder_sent_at DATETIME,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
				);
			`,
		},
		{
			name: "create_contact_notes_table",
			sql: `
				CREATE TABLE IF NOT EXISTS contact_notes (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					submission_id INTEGER NOT NULL,
					user_id INTEGER NOT NULL,
					note TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (submission_id) REFERENCES contact_submissions(id) ON DELETE CASCADE,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);
			`,
		},
		{
			name: "create_static_content_updated_at_index",
			sql: `
				CREATE INDEX IF NOT EXISTS idx_static_content_updated_at 
				ON static_content(updated_at);
			`,
		},
		{
			name: "add_needs_password_update_to_users",
			sql: `
				ALTER TABLE users ADD COLUMN needs_password_update INTEGER DEFAULT 0;
			`,
		},
		{
			name: "add_type_to_geocaches",
			sql: `
				ALTER TABLE geocaches ADD COLUMN type TEXT DEFAULT 'traditional';
			`,
		},
		{
			name: "add_placed_date_to_geocaches",
			sql: `
				ALTER TABLE geocaches ADD COLUMN placed_date DATE;
			`,
		},
		{
			name: "create_golden_key_settings_table",
			sql: `
				CREATE TABLE IF NOT EXISTS golden_key_settings (
					id INTEGER PRIMARY KEY CHECK (id = 1),
					activation_time TEXT NOT NULL DEFAULT '2026-04-12 10:12:00',
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "seed_golden_key_settings",
			sql: `
				INSERT OR IGNORE INTO golden_key_settings (id, activation_time)
				VALUES (1, '2026-04-12 10:12:00');
			`,
		},
		{
			name: "add_banner_text_to_golden_key_settings",
			sql: `
				ALTER TABLE golden_key_settings ADD COLUMN banner_text TEXT NOT NULL DEFAULT '';
			`,
		},
		{
			name: "add_rules_to_golden_key_settings",
			sql: `
				ALTER TABLE golden_key_settings ADD COLUMN rules TEXT NOT NULL DEFAULT '{}';
			`,
		},
		{
			name: "create_golden_key_months_table",
			sql: `
				CREATE TABLE IF NOT EXISTS golden_key_months (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					month_number INTEGER NOT NULL UNIQUE,
					month_name TEXT NOT NULL,
					live_date DATETIME NOT NULL,
					is_found INTEGER NOT NULL DEFAULT 0,
					finder_name TEXT,
					finder_image TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			name: "create_golden_key_hints_table",
			sql: `
				CREATE TABLE IF NOT EXISTS golden_key_hints (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					month_id INTEGER NOT NULL,
					sort_order INTEGER NOT NULL DEFAULT 0,
					content TEXT NOT NULL DEFAULT '',
					image_url TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (month_id) REFERENCES golden_key_months(id) ON DELETE CASCADE
				);
			`,
		},
		{
			name: "add_found_date_to_golden_key_months",
			sql:  `ALTER TABLE golden_key_months ADD COLUMN found_date DATETIME`,
		},
		{
			name: "seed_golden_key_months",
			sql: `
				INSERT OR IGNORE INTO golden_key_months (month_number, month_name, live_date) VALUES
				(1,  'April',     '2026-04-12 10:12:00'),
				(2,  'Mei',       '2026-05-12 10:12:00'),
				(3,  'Juni',      '2026-06-12 10:12:00'),
				(4,  'Juli',      '2026-07-12 10:12:00'),
				(5,  'Augustus',  '2026-08-12 10:12:00'),
				(6,  'September', '2026-09-12 10:12:00'),
				(7,  'Oktober',   '2026-10-12 10:12:00'),
				(8,  'November',  '2026-11-12 10:12:00'),
				(9,  'December',  '2026-12-12 10:12:00'),
				(10, 'Januari',   '2027-01-12 10:12:00'),
				(11, 'Februari',  '2027-02-12 10:12:00'),
				(12, 'Maart',     '2027-03-12 10:12:00');
			`,
		},
	}

	for _, m := range migrations {
		if _, err := db.Exec(m.sql); err != nil {
			// Ignore "duplicate column" errors for ALTER TABLE migrations
			if m.name == "add_needs_password_update_to_users" || m.name == "add_type_to_geocaches" || m.name == "add_placed_date_to_geocaches" || m.name == "add_banner_text_to_golden_key_settings" || m.name == "add_found_date_to_golden_key_months" || m.name == "add_rules_to_golden_key_settings" {
				log.Printf("  ✓ %s (column already exists)", m.name)
				continue
			}
			return fmt.Errorf("migration %s failed: %w", m.name, err)
		}
		log.Printf("  ✓ %s", m.name)
	}

	log.Println("Migrations complete!")
	return nil
}
