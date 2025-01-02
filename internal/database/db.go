package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "blocked_domains.db"

func getAdDomains() []string {
	return []string{
		"aax.amazon-adsystem.com",
		"ade.googlesyndication.com",
		"adnxs.com",
		"ads.stickyadstv.com",
		"analytics.google.com",
		"ad-delivery.net",
		"awin1.com",
		"ads.com",
		"adroll.com",
		"amplitude.nav.no",
		"appnexus.com",
		"banners.affiliatefuture.com",
		"bidswitch.net",
		"bttrack.com",
		"bingads.microsoft.com",
		"criteo.net",
		"crazyegg.com",
		"collect.mopinion.com",
		"cdn1.vntsm.com",
		"clickbank.net",
		"commissionjunction.com",
		"doubleclick.net",
		"disqus.com",
		"deploy.mopinion.com",
		"dynalyst-sync.adtdp.com",
		"f.vimeocdn.com",
		"googleads.g.doubleclick.net",
		"googleadservices.com",
		"google-analytics.com",
		"googlesyndication.com",
		"go1.aniview.com",
		"hotjar.com",
		"livefyre.com",
		"mixpanel.com",
		"mgid.com",
		"media.adfrontiers.com",
		"openx.net",
		"outbrain.com",
		"optimizely.com",
		"pubmatic.com",
		"pagead1.googlesyndication.com",
		"perfectaudience.com",
		"pagead2.googlesyndication.com",
		"rubiconproject.com",
		"revcontent.com",
		"retargeter.com",
		"refersion.com",
		"s.amazon-adsystem.com",
		"static.a-ads.com",
		"shb.richaudience.com",
		"segment.io",
		"steelhousemedia.com",
		"snowplowanalytics.com",
		"shareasale.com",
		"spot.im",
		"stats.bannersnack.com",
		"s2s.aniview.com",
		"tpe.googlesyndication.com",
		"track.wargaming-aff.com",
		"track1.aniview.com",
		"track4.aniview.com",
		"taboola.com",
		"tracker.cbx-rtb.com",
		"www.googletagmanager.com",
		"www.adsurve.com",
	}
}

/*
I guess we can toggle css styles to remove the ad placeholder as well?

common classes and selectors I found from inspecting various sites:
adsbygoogle, adsbygoogle-noablate, data-adsbygoogle-status="done, script[src*='adnxs.com'], etc

PS: I Still can't figure out how to block 'World of Tank' Ads they suck!!
*/

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS blocked_domains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL UNIQUE
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

// adds a domain to the blocked list in the database
func AddBlockedDomain(db *sql.DB, domain string) error {
	insertQuery := "INSERT OR IGNORE INTO blocked_domains (domain) VALUES (?);"
	_, err := db.Exec(insertQuery, domain)
	if err != nil {
		return fmt.Errorf("failed to add domain: %v", err)
	}
	log.Printf("Domain blocked: %s", domain)
	return nil
}

// removes a domain from the blocked list in the database
func RemoveBlockedDomain(db *sql.DB, domain string) error {
	deleteQuery := "DELETE FROM blocked_domains WHERE domain = ?;"
	_, err := db.Exec(deleteQuery, domain)
	if err != nil {
		return fmt.Errorf("failed to remove domain: %v", err)
	}
	log.Printf("Domain unblocked: %s", domain)
	return nil
}

// retrieves all blocked domains from the database
func GetBlockedDomains(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query("SELECT domain FROM blocked_domains;")
	if err != nil {
		return nil, fmt.Errorf("failed to query blocked domains: %v", err)
	}
	defer rows.Close()

	domains := make(map[string]struct{})

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, fmt.Errorf("failed to fetch a blocked domain: %v", err)
		}
		domains[domain] = struct{}{}
	}

	// simply appending the ad domains along with blocked domains map
	for _, domain := range getAdDomains() {
		domains[domain] = struct{}{}
	}

	return domains, nil
}
