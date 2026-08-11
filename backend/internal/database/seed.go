package database

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedDefaults seeds default languages, static content, and admin user if missing
func (db *DB) SeedDefaults() error {
	// Always ensure all 4 default languages exist
	log.Println("Ensuring default languages exist...")
	if err := db.seedLanguages(); err != nil {
		return err
	}

	// Check if static content exists
	var staticCount int
	db.QueryRow("SELECT COUNT(*) FROM static_content").Scan(&staticCount)

	if staticCount == 0 {
		log.Println("Seeding default static content...")
		if err := db.seedStaticContent(); err != nil {
			return err
		}
	}

	// Create default admin user if no users exist
	if err := db.seedDefaultAdmin(); err != nil {
		return err
	}

	return nil
}

// seedDefaultAdmin creates the default admin user with a random one-time password
func (db *DB) seedDefaultAdmin() error {
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)

	if userCount > 0 {
		return nil // Admin already exists
	}

	log.Println("Creating default admin user...")

	// Generate a random 32-character password
	passwordBytes := make([]byte, 16)
	if _, err := rand.Read(passwordBytes); err != nil {
		log.Printf("Warning: Failed to generate random password: %v", err)
		return err
	}
	tempPassword := hex.EncodeToString(passwordBytes)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash password: %v", err)
		return err
	}

	// Insert the default admin user with needs_password_update = 1
	_, err = db.Exec(
		`INSERT INTO users (name, email, password_hash, needs_password_update) 
		 VALUES (?, ?, ?, 1)`,
		"Admin", "admin", string(hashedPassword),
	)
	if err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
		return err
	}

	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("  DEFAULT ADMIN ACCOUNT CREATED")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("  Email:    admin")
	log.Printf("  Password: %s", tempPassword)
	log.Println("")
	log.Println("  WARNING: You MUST change this password after first login!")
	log.Println("  This password will not be shown again.")
	log.Println("═══════════════════════════════════════════════════════════")

	return nil
}

func (db *DB) seedLanguages() error {
	languages := []struct {
		code    string
		name    string
		flagURL string
	}{
		{"EN", "English", "/assets/media/fallbackLangFlags/EN.svg"},
		{"NL", "Nederlands", "/assets/media/fallbackLangFlags/NL.svg"},
		{"FR", "Français", "/assets/media/fallbackLangFlags/FR.svg"},
		{"DE", "Deutsch", "/assets/media/fallbackLangFlags/DE.svg"},
	}

	for _, lang := range languages {
		// Use INSERT OR IGNORE to only insert if not exists
		result, err := db.Exec(
			"INSERT OR IGNORE INTO languages (code, name, flag_url, active) VALUES (?, ?, ?, 1)",
			lang.code, lang.name, lang.flagURL,
		)
		if err != nil {
			log.Printf("Warning: Could not insert language %s: %v", lang.code, err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				log.Printf("  ✓ Language: %s", lang.code)
			}
		}
	}

	return nil
}

func (db *DB) seedStaticContent() error {
	// All static content with translations in EN, NL, FR, DE
	staticContent := map[string]map[string]string{
		"ButtonBack": {
			"EN": "back",
			"NL": "terug",
			"FR": "retour",
			"DE": "zurück",
		},
		"ButtonCancel": {
			"EN": "cancel",
			"NL": "annuleren",
			"FR": "annuler",
			"DE": "abbrechen",
		},
		"ButtonConfirm": {
			"EN": "confirm",
			"NL": "bevestigen",
			"FR": "confirmer",
			"DE": "bestätigen",
		},
		"ButtonNext": {
			"EN": "next",
			"NL": "volgende",
			"FR": "suivant",
			"DE": "weiter",
		},
		"ButtonPurchase": {
			"EN": "purchase ticket(s)",
			"NL": "ticket(s) kopen",
			"FR": "acheter le(s) billet(s)",
			"DE": "ticket(s) kaufen",
		},
		"ContactCallTxt": {
			"EN": "give us a call",
			"NL": "bel ons",
			"FR": "appelez-nous",
			"DE": "rufen Sie uns an",
		},
		"ContactHelpQuestionTxt": {
			"EN": "can we help you?",
			"NL": "kunnen wij u helpen?",
			"FR": "pouvons-nous vous aider?",
			"DE": "können wir Ihnen helfen?",
		},
		"ContactHelpTxt": {
			"EN": "Do you have a question? Is something bothering you, or do you simply want to tell us something? Then you can do so through this channel; Mail, write, call or fill out the form.",
			"NL": "Heeft u een vraag? Zit er iets dwars, of wilt u ons gewoon iets vertellen? Dan kan dat via dit kanaal; Mail, schrijf, bel of vul het formulier in.",
			"FR": "Avez-vous une question? Quelque chose vous dérange ou vous souhaitez simplement nous dire quelque chose? Vous pouvez le faire par ce canal; Envoyez un mail, écrivez, appelez ou remplissez le formulaire.",
			"DE": "Haben Sie eine Frage? Stört Sie etwas, oder möchten Sie uns einfach etwas mitteilen? Dann können Sie dies über diesen Kanal tun; Mailen, schreiben, anrufen oder das Formular ausfüllen.",
		},
		"ContactMailTxt": {
			"EN": "send us an e-mail",
			"NL": "stuur ons een e-mail",
			"FR": "envoyez-nous un e-mail",
			"DE": "senden Sie uns eine E-Mail",
		},
		"ContactPostalMailTxt": {
			"EN": "postal mail",
			"NL": "post",
			"FR": "courrier postal",
			"DE": "Post",
		},
		"FormAddress": {
			"EN": "address",
			"NL": "adres",
			"FR": "adresse",
			"DE": "Adresse",
		},
		"FormFirstName": {
			"EN": "firstname",
			"NL": "voornaam",
			"FR": "prénom",
			"DE": "Vorname",
		},
		"FormFullName": {
			"EN": "fullname",
			"NL": "volledige naam",
			"FR": "nom complet",
			"DE": "vollständiger Name",
		},
		"FormLastName": {
			"EN": "lastname",
			"NL": "achternaam",
			"FR": "nom de famille",
			"DE": "Nachname",
		},
		"FormMail": {
			"EN": "e-mail",
			"NL": "e-mail",
			"FR": "e-mail",
			"DE": "E-Mail",
		},
		"FormMessage": {
			"EN": "message",
			"NL": "bericht",
			"FR": "message",
			"DE": "Nachricht",
		},
		"FormPhoneNumber": {
			"EN": "phone number",
			"NL": "telefoonnummer",
			"FR": "numéro de téléphone",
			"DE": "Telefonnummer",
		},
		"FormSearch": {
			"EN": "search",
			"NL": "zoeken",
			"FR": "rechercher",
			"DE": "suchen",
		},
		"FormSubject": {
			"EN": "subject",
			"NL": "onderwerp",
			"FR": "sujet",
			"DE": "Betreff",
		},
		"FormSubmit": {
			"EN": "submit",
			"NL": "verzenden",
			"FR": "envoyer",
			"DE": "einreichen",
		},
		"FormSuccess": {
			"EN": "sent",
			"NL": "verzonden",
			"FR": "envoyé",
			"DE": "gesendet",
		},
		"FormFailed": {
			"EN": "something went wrong, couldn't submit",
			"NL": "er liep iets mis, kon formulier niet verzenden",
			"FR": "un problème est survenu, le formulaire n'a pas pu être envoyé",
			"DE": "etwas ist schief gelaufen, das Formular konnte nicht gesendet werden",
		},
		"NavEvents": {
			"EN": "events",
			"NL": "evenementen",
			"FR": "événements",
			"DE": "Veranstaltungen",
		},
		"NavGeocaches": {
			"EN": "geocaches",
			"NL": "geocaches",
			"FR": "géocaches",
			"DE": "Geocaches",
		},
		"NavHome": {
			"EN": "home",
			"NL": "home",
			"FR": "accueil",
			"DE": "Startseite",
		},
		"NavShop": {
			"EN": "shop",
			"NL": "winkel",
			"FR": "boutique",
			"DE": "Shop",
		},
		"SocialsFollowTxt": {
			"EN": "follow us on social media",
			"NL": "volg ons op sociale media",
			"FR": "suivez-nous sur les réseaux sociaux",
			"DE": "folgen Sie uns in den sozialen Medien",
		},
		"SplashBody": {
			"EN": "Every year, we organise a range of Geocaching events that bring people together. These events include activities such as our clean-up effort, known as CITO (Cache In Trash Out), the yearly GIFF-Filmfestival showcasing films made by geocachers worldwide, and our flagship geocaching event, the Brugse Beer.",
			"NL": "Jaarlijks organiseren wij diverse geocaching evenementen waarbij we mensen samenbrengen. Deze evenementen omvatten onder andere onze zwerfvuilactie, CITO (Cache In Trash Out), het jaarlijkse GIFF-filmfestival, dat gewijd is aan films gemaakt door geocachers van over de hele wereld, en ons vooraanstaande geocaching evenement, de BRUGSE BEER.",
			"FR": "Chaque année, nous organisons une variété d'événements de géocaching qui rassemblent les gens. Ces événements comprennent des activités telles que notre effort de nettoyage, connu sous le nom de CITO (Cache In Trash Out), le festival annuel GIFF qui présente des films réalisés par des géocacheurs du monde entier, ainsi que notre événement phare de géocaching, le Brugse Beer.",
			"DE": "Jedes Jahr organisieren wir eine Vielzahl von Geocaching-Veranstaltungen, die Menschen zusammenbringen. Diese Veranstaltungen umfassen Aktivitäten wie unsere Aufräumaktion, bekannt als CITO (Cache In Trash Out), das jährliche GIFF-Filmfestival, bei dem Filme von Geocachern aus der ganzen Welt gezeigt werden, und unser Hauptgeocaching-Event, der Brugse Beer",
		},
		"SplashImg": {
			"EN": "static/bertje.jpg",
			"NL": "static/bertje.jpg",
			"FR": "static/bertje.jpg",
			"DE": "static/bertje.jpg",
		},
		"SplashTitle": {
			"EN": "GeocachingBrughia brings people together",
			"NL": "GeocachingBrughia brengt mensen samen",
			"FR": "GeocachingBrughia rassemble les gens",
			"DE": "GeocachingBrughia vereint Menschen",
		},
		"UILoadingEvents": {
			"EN": "loading events",
			"NL": "evenementen laden",
			"FR": "chargement des événements",
			"DE": "Veranstaltungen werden geladen",
		},
		"UILoadingGeocaches": {
			"EN": "loading geocaches",
			"NL": "geocaches laden",
			"FR": "chargement des géocaches",
			"DE": "Geocaches werden geladen",
		},
		"UILoadingSite": {
			"EN": "loading website content",
			"NL": "website-inhoud laden",
			"FR": "chargement du contenu du site",
			"DE": "Website-Inhalt wird geladen",
		},
		"UINoEvents": {
			"EN": "There are no events planned at the moment.",
			"NL": "Er zijn momenteel geen evenementen gepland.",
			"FR": "Il n'y a pas d'événements prévus pour le moment.",
			"DE": "Momentan sind keine Veranstaltungen geplant.",
		},
		"UINoEventsSubTxt": {
			"EN": "Check back later!",
			"NL": "Kijk later nog eens!",
			"FR": "Revenez plus tard!",
			"DE": "Schauen Sie später noch einmal vorbei!",
		},
		"UINoGeocaches": {
			"EN": "There are no geocaches at the moment.",
			"NL": "Er zijn momenteel geen geocaches.",
			"FR": "Il n'y a pas de géocaches pour le moment.",
			"DE": "Momentan gibt es keine Geocaches.",
		},
		"UINoGeocachesSubTxt": {
			"EN": "Check back later!",
			"NL": "Kijk later nog eens!",
			"FR": "Revenez plus tard!",
			"DE": "Schauen Sie später noch einmal vorbei!",
		},
		"UINoStoreItems": {
			"EN": "Items not available.",
			"NL": "Artikelen niet beschikbaar.",
			"FR": "Articles non disponibles.",
			"DE": "Artikel nicht verfügbar.",
		},
		"UINoStoreItemsSubTxt": {
			"EN": "The store item information could not be retrieved.",
			"NL": "De informatie over de artikelen kon niet worden opgehaald.",
			"FR": "Les informations sur les articles n'ont pas pu être récupérées.",
			"DE": "Die Artikelinformationen konnten nicht abgerufen werden.",
		},
		"UIShopEmpty": {
			"EN": "The shop is currently empty.",
			"NL": "De webshop is momenteel leeg.",
			"FR": "La boutique est actuellement vide.",
			"DE": "Der Shop ist derzeit leer.",
		},
		"UIShopEmptySubTxt": {
			"EN": "Check back later!",
			"NL": "Kom later terug!",
			"FR": "Revenez plus tard !",
			"DE": "Schauen Sie später noch einmal vorbei!",
		},
		"UILoadingShop": {
			"EN": "Loading shop",
			"NL": "Winkel laden",
			"FR": "Chargement de la boutique",
			"DE": "Shop wird geladen",
		},
		"ShopTicketsTitle": {
			"EN": "Tickets",
			"NL": "Tickets",
			"FR": "Billets",
			"DE": "Tickets",
		},
		"ShopTicketsSubTxt": {
			"EN": "Order your tickets for our events.",
			"NL": "Bestel je tickets voor onze evenementen.",
			"FR": "Commandez vos billets pour nos événements.",
			"DE": "Bestellen Sie Ihre Tickets für unsere Veranstaltungen.",
		},
		"ShopTicketsGoTo": {
			"EN": "Go to the ticket shop",
			"NL": "Ga naar de ticketshop",
			"FR": "Aller à la billetterie",
			"DE": "Zur Ticket-Shop",
		},
		"ShopMerchTitle": {
			"EN": "Shop",
			"NL": "Webshop",
			"FR": "Boutique",
			"DE": "Shop",
		},
		"ShopMerchSubTxt": {
			"EN": "Merchandise and other items.",
			"NL": "Merchandise en andere items.",
			"FR": "Produits dérivés et autres articles.",
			"DE": "Merchandise und andere Artikel.",
		},
		"ShopInStock": {
			"EN": "In stock",
			"NL": "Op voorraad",
			"FR": "En stock",
			"DE": "Auf Lager",
		},
		"ShopStockCount": {
			"EN": "in stock",
			"NL": "op voorraad",
			"FR": "en stock",
			"DE": "auf Lager",
		},
		"ShopPickup": {
			"EN": "Pickup",
			"NL": "Afhalen",
			"FR": "Retrait",
			"DE": "Abholung",
		},
		"ShopShipping": {
			"EN": "Shipping",
			"NL": "Verzenden",
			"FR": "Expédition",
			"DE": "Versand",
		},
		"ShopBuy": {
			"EN": "Buy",
			"NL": "Kopen",
			"FR": "Acheter",
			"DE": "Kaufen",
		},
		"ShopCountries": {
			"EN": "countries",
			"NL": "landen",
			"FR": "pays",
			"DE": "Länder",
		},
		"ShopOrderTitle": {
			"EN": "Order",
			"NL": "Bestellen",
			"FR": "Commander",
			"DE": "Bestellen",
		},
		"ShopEmail": {
			"EN": "Email",
			"NL": "E-mailadres",
			"FR": "E-mail",
			"DE": "E-Mail",
		},
		"ShopQuantity": {
			"EN": "Quantity",
			"NL": "Aantal",
			"FR": "Quantité",
			"DE": "Anzahl",
		},
		"ShopDelivery": {
			"EN": "Delivery",
			"NL": "Levering",
			"FR": "Livraison",
			"DE": "Lieferung",
		},
		"ShopName": {
			"EN": "Name",
			"NL": "Naam",
			"FR": "Nom",
			"DE": "Name",
		},
		"ShopAddress": {
			"EN": "Address",
			"NL": "Adres",
			"FR": "Adresse",
			"DE": "Adresse",
		},
		"ShopCity": {
			"EN": "City",
			"NL": "Stad",
			"FR": "Ville",
			"DE": "Stadt",
		},
		"ShopPostalCode": {
			"EN": "Postal code",
			"NL": "Postcode",
			"FR": "Code postal",
			"DE": "Postleitzahl",
		},
		"ShopCountry": {
			"EN": "Country",
			"NL": "Land",
			"FR": "Pays",
			"DE": "Land",
		},
		"ShopTotal": {
			"EN": "Total:",
			"NL": "Totaal:",
			"FR": "Total :",
			"DE": "Gesamt:",
		},
		"ShopStripeNote": {
			"EN": "You will be redirected to the secure Stripe payment environment.",
			"NL": "Je wordt doorgestuurd naar de beveiligde betaalomgeving van Stripe.",
			"FR": "Vous serez redirigé vers l'environnement de paiement sécurisé de Stripe.",
			"DE": "Sie werden zur sicheren Zahlungsumgebung von Stripe weitergeleitet.",
		},
		"ShopPay": {
			"EN": "Pay",
			"NL": "Betalen",
			"FR": "Payer",
			"DE": "Bezahlen",
		},
		"ShopProcessing": {
			"EN": "Processing...",
			"NL": "Bezig...",
			"FR": "Traitement en cours...",
			"DE": "Wird verarbeitet...",
		},
		"ShopClose": {
			"EN": "Close",
			"NL": "Sluiten",
			"FR": "Fermer",
			"DE": "Schließen",
		},
		"ShopLoadError": {
			"EN": "Could not load the shop.",
			"NL": "Kon de winkel niet laden.",
			"FR": "Impossible de charger la boutique.",
			"DE": "Der Shop konnte nicht geladen werden.",
		},
		"ShopErrorEmail": {
			"EN": "Email is required",
			"NL": "E-mailadres is verplicht",
			"FR": "L'e-mail est obligatoire",
			"DE": "E-Mail ist erforderlich",
		},
		"ShopErrorQuantity": {
			"EN": "Quantity must be at least 1",
			"NL": "Aantal moet minimaal 1 zijn",
			"FR": "La quantité doit être au moins 1",
			"DE": "Die Anzahl muss mindestens 1 sein",
		},
		"ShopErrorShipping": {
			"EN": "All shipping fields are required",
			"NL": "Alle verzendvelden zijn verplicht",
			"FR": "Tous les champs d'expédition sont obligatoires",
			"DE": "Alle Versandfelder sind erforderlich",
		},
		"ShopErrorStart": {
			"EN": "Failed to start payment.",
			"NL": "Betaling starten mislukt.",
			"FR": "Échec du démarrage du paiement.",
			"DE": "Zahlung konnte nicht gestartet werden.",
		},
		"ShopErrorGeneric": {
			"EN": "An error occurred.",
			"NL": "Er is een fout opgetreden.",
			"FR": "Une erreur s'est produite.",
			"DE": "Ein Fehler ist aufgetreten.",
		},
		"CountryBE": {
			"EN": "Belgium",
			"NL": "België",
			"FR": "Belgique",
			"DE": "Belgien",
		},
		"CountryNL": {
			"EN": "Netherlands",
			"NL": "Nederland",
			"FR": "Pays-Bas",
			"DE": "Niederlande",
		},
		"CountryLU": {
			"EN": "Luxembourg",
			"NL": "Luxemburg",
			"FR": "Luxembourg",
			"DE": "Luxemburg",
		},
		"CountryFR": {
			"EN": "France",
			"NL": "Frankrijk",
			"FR": "France",
			"DE": "Frankreich",
		},
		"CountryDE": {
			"EN": "Germany",
			"NL": "Duitsland",
			"FR": "Allemagne",
			"DE": "Deutschland",
		},
		"CountryGB": {
			"EN": "United Kingdom",
			"NL": "Verenigd Koninkrijk",
			"FR": "Royaume-Uni",
			"DE": "Vereinigtes Königreich",
		},
		"CountryIE": {
			"EN": "Ireland",
			"NL": "Ierland",
			"FR": "Irlande",
			"DE": "Irland",
		},
		"CountryES": {
			"EN": "Spain",
			"NL": "Spanje",
			"FR": "Espagne",
			"DE": "Spanien",
		},
		"CountryPT": {
			"EN": "Portugal",
			"NL": "Portugal",
			"FR": "Portugal",
			"DE": "Portugal",
		},
		"CountryIT": {
			"EN": "Italy",
			"NL": "Italië",
			"FR": "Italie",
			"DE": "Italien",
		},
		"CountryAT": {
			"EN": "Austria",
			"NL": "Oostenrijk",
			"FR": "Autriche",
			"DE": "Österreich",
		},
		"CountryCH": {
			"EN": "Switzerland",
			"NL": "Zwitserland",
			"FR": "Suisse",
			"DE": "Schweiz",
		},
		"CountryDK": {
			"EN": "Denmark",
			"NL": "Denemarken",
			"FR": "Danemark",
			"DE": "Dänemark",
		},
		"CountrySE": {
			"EN": "Sweden",
			"NL": "Zweden",
			"FR": "Suède",
			"DE": "Schweden",
		},
		"CountryNO": {
			"EN": "Norway",
			"NL": "Noorwegen",
			"FR": "Norvège",
			"DE": "Norwegen",
		},
		"CountryFI": {
			"EN": "Finland",
			"NL": "Finland",
			"FR": "Finlande",
			"DE": "Finnland",
		},
		"CountryPL": {
			"EN": "Poland",
			"NL": "Polen",
			"FR": "Pologne",
			"DE": "Polen",
		},
		"CountryCZ": {
			"EN": "Czechia",
			"NL": "Tsjechië",
			"FR": "Tchéquie",
			"DE": "Tschechien",
		},
		"CountrySK": {
			"EN": "Slovakia",
			"NL": "Slowakije",
			"FR": "Slovaquie",
			"DE": "Slowakei",
		},
		"CountryHU": {
			"EN": "Hungary",
			"NL": "Hongarije",
			"FR": "Hongrie",
			"DE": "Ungarn",
		},
		"CountryRO": {
			"EN": "Romania",
			"NL": "Roemenië",
			"FR": "Roumanie",
			"DE": "Rumänien",
		},
		"CountryBG": {
			"EN": "Bulgaria",
			"NL": "Bulgarije",
			"FR": "Bulgarie",
			"DE": "Bulgarien",
		},
		"CountryHR": {
			"EN": "Croatia",
			"NL": "Kroatië",
			"FR": "Croatie",
			"DE": "Kroatien",
		},
		"CountrySI": {
			"EN": "Slovenia",
			"NL": "Slovenië",
			"FR": "Slovénie",
			"DE": "Slowenien",
		},
		"CountryEE": {
			"EN": "Estonia",
			"NL": "Estland",
			"FR": "Estonie",
			"DE": "Estland",
		},
		"CountryLV": {
			"EN": "Latvia",
			"NL": "Letland",
			"FR": "Lettonie",
			"DE": "Lettland",
		},
		"CountryLT": {
			"EN": "Lithuania",
			"NL": "Litouwen",
			"FR": "Lituanie",
			"DE": "Litauen",
		},
		"CountryGR": {
			"EN": "Greece",
			"NL": "Griekenland",
			"FR": "Grèce",
			"DE": "Griechenland",
		},
		"UIPageNotFound": {
			"EN": "Page Not Found",
			"NL": "Pagina niet gevonden",
			"FR": "Page non trouvée",
			"DE": "Seite nicht gefunden",
		},
		"UIPageNotFoundSubTxt": {
			"EN": "Oops! The page you are looking for could not be found. You will be redirected back to the homepage in: ///t/// seconds.",
			"NL": "Oeps! De pagina die u zoekt, kon niet worden gevonden. U wordt binnen ///t/// seconden teruggeleid naar de homepage.",
			"FR": "Oops ! La page que vous recherchez est introuvable. Vous serez redirigé vers la page d'accueil dans ///t/// secondes.",
			"DE": "Oops! Die gesuchte Seite konnte nicht gefunden werden. Sie werden in ///t/// Sekunden zur Startseite weitergeleitet.",
		},
		"GeocacheTitle": {
			"EN": "title",
			"NL": "titel",
			"FR": "titre",
			"DE": "titel",
		},
		"GeocacheGeoLink": {
			"EN": "link to cache",
			"NL": "link naar cache",
			"FR": "lien vers la cache",
			"DE": "Link zur Cache",
		},
		"GeocacheDifficulty": {
			"EN": "difficulty",
			"NL": "moeilijkheidsgraad",
			"FR": "difficulté",
			"DE": "Schwierigkeit",
		},
		"GeocacheTerrain": {
			"EN": "terrain",
			"NL": "terrein",
			"FR": "terrain",
			"DE": "Gelände",
		},
		// Contact information (same for all languages since they are data, not translations)
		"ContactAddress": {
			"EN": "Korte Kwadeplasstraat 6, 8020 Oostkamp",
			"NL": "Korte Kwadeplasstraat 6, 8020 Oostkamp",
			"FR": "Korte Kwadeplasstraat 6, 8020 Oostkamp",
			"DE": "Korte Kwadeplasstraat 6, 8020 Oostkamp",
		},
		"ContactAddressUrl": {
			"EN": "https://www.google.com/maps/place/Korte+Kwadeplasstraat+6,+8020+Oostkamp",
			"NL": "https://www.google.com/maps/place/Korte+Kwadeplasstraat+6,+8020+Oostkamp",
			"FR": "https://www.google.com/maps/place/Korte+Kwadeplasstraat+6,+8020+Oostkamp",
			"DE": "https://www.google.com/maps/place/Korte+Kwadeplasstraat+6,+8020+Oostkamp",
		},
		"ContactPhone1": {
			"EN": "+32 50 841 331",
			"NL": "+32 50 841 331",
			"FR": "+32 50 841 331",
			"DE": "+32 50 841 331",
		},
		"ContactPhone2": {
			"EN": "+32 487 906 431",
			"NL": "+32 487 906 431",
			"FR": "+32 487 906 431",
			"DE": "+32 487 906 431",
		},
		"ContactEmail": {
			"EN": "info@geocachingbrughia.be",
			"NL": "info@geocachingbrughia.be",
			"FR": "info@geocachingbrughia.be",
			"DE": "info@geocachingbrughia.be",
		},
	}

	for property, translations := range staticContent {
		for langCode, content := range translations {
			_, err := db.Exec(
				"INSERT INTO static_content (property, lang_code, content) VALUES (?, ?, ?)",
				property, langCode, content,
			)
			if err != nil {
				log.Printf("Warning: Could not insert static content %s/%s: %v", property, langCode, err)
			}
		}
	}

	log.Printf("  ✓ Seeded %d translation keys", len(staticContent))
	return nil
}
