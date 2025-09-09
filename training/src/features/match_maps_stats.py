import pandas as pd


def compute_match_map_stats(players_stats_df: pd.DataFrame, opp_players_stats_df: pd.DataFrame):
    def f(row: pd.Series):
        map_players_stats_df = players_stats_df[(players_stats_df["map_id"] == row["map_id"]) & (players_stats_df["side"] == "def")]
        opp_map_players_stats_df = opp_players_stats_df[(opp_players_stats_df["map_id"] == row["map_id"]) & (players_stats_df["side"] == "def")]
        team_id = map_players_stats_df.iloc[0]["team_id"]

        team_comps = sorted(map_players_stats_df["role"].to_list())
        team_against_comps = sorted(opp_map_players_stats_df["role"].to_list())
        print(team_comps)

        return pd.Series(
            {
                "team_def_score": row["team_1_def_score"] if row["team_1_id"] == team_id else row["team_2_def_score"],
                "team_atk_score": row["team_1_atk_score"] if row["team_1_id"] == team_id else row["team_2_atk_score"],
                "team_ot_score": row["team_1_ot_score"] if row["team_1_id"] == team_id else row["team_2_ot_score"],
                "team_against_def_score": row["team_1_def_score"] if row["team_1_id"] != team_id else row["team_2_def_score"],
                "team_against_atk_score": row["team_1_atk_score"] if row["team_1_id"] != team_id else row["team_2_atk_score"],
                "team_against_ot_score": row["team_1_ot_score"] if row["team_1_id"] != team_id else row["team_2_ot_score"],
                "duration": row["duration"],
                "team_agents": sorted(map_players_stats_df["agent_id"].unique().tolist()),
                "team_against_agents": sorted(opp_map_players_stats_df["agent_id"].unique().tolist()),
                "team_duelists": team_comps.count("duelist"),
                "team_controllers": team_comps.count("controller"),
                "team_sentinels": team_comps.count("sentinel"),
                "team_initiators": team_comps.count("initiator"),
                "team_against_duelists": team_against_comps.count("duelist"),
                "team_against_controllers": team_against_comps.count("controller"),
                "team_against_sentinels": team_against_comps.count("sentinel"),
                "team_against_initiators": team_against_comps.count("initiator"),
            }
        )

    return f


def get_match_maps_stats(match_maps_df: pd.DataFrame, players_stats_df: pd.DataFrame, opp_players_stats_df: pd.DataFrame):
    return match_maps_df.apply(compute_match_map_stats(players_stats_df, opp_players_stats_df), axis=1)
