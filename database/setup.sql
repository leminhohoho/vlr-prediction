DROP TABLE IF EXISTS matches;
CREATE TABLE IF NOT EXISTS matches (
    id INTEGER PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    date TEXT NOT NULL,
    tournament_id INTEGER NOT NULL,
    stage TEXT CHECK(stage IN ('group_stage', 'playoff', 'grand_final')),
    team_1_id INTEGER NOT NULL,
    team_2_id INTEGER NOT NULL,
    team_1_score INTEGER NOT NULL CHECK(team_1_score >= 0),
    team_2_score INTEGER NOT NULL CHECK(team_2_score >= 0),
    team_1_rating INTEGER CHECK(team_1_rating >= 0),
    team_2_rating INTEGER CHECK(team_2_rating >= 0),

    FOREIGN KEY (tournament_id) REFERENCES tournaments(id),
    FOREIGN KEY (team_1_id) REFERENCES teams(id),
    FOREIGN KEY (team_2_id) REFERENCES teams(id)
);

DROP TABLE IF EXISTS match_maps;
CREATE TABLE IF NOT EXISTS matche_maps (
    match_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    duration INTEGER CHECK (duration >= 0),
    team_1_id INTEGER NOT NULL,
    team_2_id INTEGER NOT NULL,
    team_1_def_score INTEGER NOT NULL CHECK(team_1_def_score >= 0),
    team_1_atk_score INTEGER NOT NULL CHECK(team_1_atk_score >= 0),
    team_1_ot_score INTEGER NOT NULL CHECK(team_1_ot_score >= 0),
    team_2_def_score INTEGER NOT NULL CHECK(team_2_def_score >= 0),
    team_2_atk_score INTEGER NOT NULL CHECK(team_2_atk_score >= 0),
    team_2_ot_score INTEGER NOT NULL CHECK(team_2_ot_score >= 0),
    team_def_first INTEGER NOT NULL,
    team_pick INTEGER NOT NULL,

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (team_1_id) REFERENCES teams(id),
    FOREIGN KEY (team_2_id) REFERENCES teams(id),
    FOREIGN KEY (team_def_first) REFERENCES teams(id),
    FOREIGN KEY (team_pick) REFERENCES teams(id)
);

DROP TABLE IF EXISTS round_stats;
CREATE TABLE IF NOT EXISTS round_stats (
    match_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    round_no INTEGER NOT NULL CHECK(round_no > 0),
    team_1_id INTEGER NOT NULL,
    team_2_id INTEGER NOT NULL,
    team_def INTEGER NOT NULL,
    team_1_buy_type TEXT NOT NULL CHECK(team_1_buy_type IN ('pistol', 'eco', 'semi_eco', 'semi_buy', 'full_buy')),
    team_2_buy_type TEXT NOT NULL CHECK(team_2_buy_type IN ('pistol', 'eco', 'semi_eco', 'semi_buy', 'full_buy')),
    team_1_bank INTEGER NOT NULL CHECK(team_1_bank >= 0),
    team_2_bank INTEGER NOT NULL CHECK(team_2_bank >= 0),
    team_won INTEGER NOT NULL,
    won_method TEXT NOT NULL CHECK(won_method IN ('eliminate', 'spike_explode', 'defuse')),

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (team_1_id) REFERENCES teams(id),
    FOREIGN KEY (team_2_id) REFERENCES teams(id),
    FOREIGN KEY (team_def) REFERENCES teams(id),
    FOREIGN KEY (team_won) REFERENCES teams(id)
);

DROP TABLE IF EXISTS players_duel_stats;
CREATE TABLE IF NOT EXISTS players_duel_stats (
    match_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    team_1_player_id INTEGER NOT NULL,
    team_2_player_id INTEGER NOT NULL,
    team_1_player_kills_vs_team_2_player INTEGER NOT NULL CHECK(team_1_player_kills_vs_team_2_player >= 0),
    team_2_player_kills_vs_team_1_player INTEGER NOT NULL CHECK(team_2_player_kills_vs_team_1_player >= 0),
    team_1_player_first_kills_vs_team_2_player INTEGER NOT NULL CHECK(team_1_player_first_kills_vs_team_2_player >= 0),
    team_2_player_first_kills_vs_team_1_player INTEGER NOT NULL CHECK(team_2_player_first_kills_vs_team_1_player >= 0),
    team_1_player_op_kills_vs_team_2_player INTEGER NOT NULL CHECK(team_1_player_op_kills_vs_team_2_player >= 0),
    team_2_player_op_kills_vs_team_1_player INTEGER NOT NULL CHECK(team_2_player_op_kills_vs_team_1_player >= 0),

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (team_1_player_id) REFERENCES players(id),
    FOREIGN KEY (team_2_player_id) REFERENCES players(id)
);

DROP TABLE IF EXISTS player_overview_stats;
CREATE TABLE IF NOT EXISTS player_overview_stats (
    match_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    side TEXT NOT NULL CHECK(side IN ('def', 'atk')),
    rating REAL NOT NULL CHECK(rating >= 0),
    acs REAL NOT NULL CHECK(acs >= 0),
    kills REAL NOT NULL,
    deaths REAL NOT NULL CHECK(deaths >= 0),
    assists REAL NOT NULL CHECK(assists >= 0),
    kast REAL NOT NULL CHECK(kast >= 0 AND kast <= 100),
    adr REAL NOT NULL CHECK(adr >= 0),
    hs REAL NOT NULL CHECK(hs >= 0 AND hs <= 100),
    first_kills REAL NOT NULL CHECK(first_kills >= 0),
    first_deaths REAL NOT NULL CHECK(first_deaths >= 0),

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (player_id) REFERENCES players(id)
);

DROP TABLE IF EXISTS player_highlights;
CREATE TABLE IF NOT EXISTS player_highlights (
    match_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    round_no INTEGER NOT NULL CHECK(round_no > 0),
    team_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('maker', 'receiver')),
    type TEXT NOT NULL CHECK(type IN ('2k', '3k', '4k', '5k', '1v1', '1v2', '1v3', '1v4', '1v5')),

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (player_id) REFERENCES players(id)
);

DROP TABLE IF EXISTS ban_pick_log;
CREATE TABLE IF NOT EXISTS ban_pick_log (
    match_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    map_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK(action IN ('ban', 'pick')),
    ban_pick_order INTEGER NOT NULL CHECK(ban_pick_order >= 0),

    FOREIGN KEY (match_id) REFERENCES matches(id),
    FOREIGN KEY (map_id) REFERENCES maps(id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

DROP TABLE IF EXISTS tournaments;
CREATE TABLE IF NOT EXISTS tournaments (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    prize_pool INTEGER NOT NULL,
    tier_1 INTEGER NOT NULL
);

DROP TABLE IF EXISTS players;
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    img_url TEXT
);

DROP TABLE IF EXISTS teams;
CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    shorthand_name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    img_url TEXT,
    region_id INTEGER NOT NULL,

    FOREIGN KEY (region_id) REFERENCES regions(id)
);

DROP TABLE IF EXISTS maps;
CREATE TABLE IF NOT EXISTS maps (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    release_date TEXT NOT NULL
);

DROP TABLE IF EXISTS maps_pool;
CREATE TABLE IF NOT EXISTS maps_pool (
    date TEXT NOT NULL,
    patch_no TEXT NOT NULL, 
    map_id INTEGER NOT NULL,

    FOREIGN KEY (map_id) REFERENCES maps(id)
);

DROP TABLE IF EXISTS regions;
CREATE TABLE IF NOT EXISTS regions (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

DROP TABLE IF EXISTS agents;
CREATE TABLE IF NOT EXISTS agents (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    agent_type TEXT NOT NULL CHECK(agent_type IN ('duelist', 'controller', 'sentinel', 'initiator')),
    release_date TEXT NOT NULL
);
-- MAPS, MAPS_POOL, REGIONS AND AGENT INSERT --
INSERT OR IGNORE INTO maps(name, release_date) VALUES
('Ascent', '2020-06-02'),
('Bind', '2020-04-07'),
('Breeze', '2021-04-27'),
('Haven', '2020-04-07'),
('Icebox', '2020-10-13'),
('Lotus', '2023-01-10'),
('Pearl', '2022-06-22'),
('Split', '2020-04-07'),
('Sunset', '2023-08-29'),
('Abyss', '2024-06-11'),
('Fracture', '2021-09-08');

INSERT OR IGNORE INTO maps_pool(date, patch_no, map_id) VALUES
-- June 2, 2020 (Patch 1.0, Game Launch)
('2020-06-02', '1.0', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2020-06-02', '1.0', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2020-06-02', '1.0', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2020-06-02', '1.0', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
-- October 13, 2020 (Patch 1.08, Episode 01, Act 3)
('2020-10-13', '1.08', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2020-10-13', '1.08', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2020-10-13', '1.08', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2020-10-13', '1.08', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2020-10-13', '1.08', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
-- April 27, 2021 (Patch 2.08, Episode 02, Act 3)
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2021-04-27', '2.08', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
-- September 8, 2021 (Patch 3.05, Episode 03, Act 2)
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2021-09-08', '3.05', (SELECT id FROM maps WHERE name = 'Fracture' LIMIT 1)),
-- June 22, 2022 (Patch 5.0, Episode 05, Act 1)
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Fracture' LIMIT 1)),
('2022-06-22', '5.0', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
-- January 10, 2023 (Patch 6.0, Episode 06, Act 1)
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Fracture' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2023-01-10', '6.0', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
-- April 25, 2023 (Patch 6.08, Episode 06, Act 3)
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2023-04-25', '6.08', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
-- August 29, 2023 (Patch 7.04, Episode 07, Act 2)
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1)),
('2023-08-29', '7.04', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
-- January 9, 2024 (Patch 8.0, Episode 08, Act 1)
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2024-01-09', '8.0', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
-- June 11, 2024 (Patch 8.11, Episode 08, Act 3)
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Breeze' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Abyss' LIMIT 1)),
('2024-06-11', '8.11', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
-- October 22, 2024 (Patch 9.08, Episode 09, Act 3)
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Abyss' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2024-10-22', '9.08', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
-- January 8, 2025 (Patch 10.00, Season 2025, Episode 10, Act 1)
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Abyss' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Bind' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Fracture' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2025-01-08', '10.00', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
-- March 4, 2025 (Patch 10.04, Season 2025, Act 2)
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Fracture' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2025-03-04', '10.04', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
-- April 29, 2025 (Patch 10.08, Season 2025, Act 3)
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2025-04-29', '10.08', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1)),
-- June 10, 2025 (Patch 10.10, Season 2025, Act 3, Estimated)
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Ascent' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Haven' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Icebox' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Lotus' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Pearl' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Split' LIMIT 1)),
('2025-06-10', '10.10', (SELECT id FROM maps WHERE name = 'Sunset' LIMIT 1));

INSERT OR IGNORE INTO agents(name, agent_type, release_date) VALUES
('Astra', 'controller', '2021-03-02'),
('Breach', 'initiator', '2020-04-07'),
('Brimstone', 'controller', '2020-04-07'),
('Chamber', 'sentinel', '2021-11-16'),
('Clove', 'controller', '2024-03-26'),
('Cypher', 'sentinel', '2020-04-07'),
('Deadlock', 'sentinel', '2023-08-29'),
('Fade', 'initiator', '2022-04-27'),
('Gekko', 'initiator', '2023-03-07'),
('Harbor', 'controller', '2022-10-18'),
('Iso', 'duelist', '2023-10-31'),
('Jett', 'duelist', '2020-04-07'),
('Kayo', 'initiator', '2021-06-22'),
('Killjoy', 'sentinel', '2020-08-04'),
('Neon', 'duelist', '2022-01-11'),
('Omen', 'controller', '2020-04-07'),
('Phoenix', 'duelist', '2020-04-07'),
('Raze', 'duelist', '2020-04-07'),
('Reyna', 'duelist', '2020-06-02'),
('Sage', 'sentinel', '2020-04-07'),
('Skye', 'initiator', '2020-10-27'),
('Sova', 'initiator', '2020-04-07'),
('Viper', 'controller', '2020-04-07'),
('Vyse', 'sentinel', '2024-08-28'),
('Yoru', 'duelist', '2021-01-12'),
('Tejo', 'initiator', '2025-01-08'),
('Waylay', 'duelist', '2025-03-04');

INSERT OR IGNORE INTO regions(name) VALUES
('Americas'),
('EMEA'),
('Pacific'),
('China'),
('International');

SELECT * FROM maps;
SELECT * FROM agents;
SELECT * FROM regions;
SELECT * FROM maps_pool ORDER BY date;

