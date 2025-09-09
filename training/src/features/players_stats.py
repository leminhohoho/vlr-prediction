import pandas as pd
import numpy as np
from src.utils import load_players_stats


def get_player_stats(enriched_highlights_df: pd.DataFrame, rounds_df: pd.DataFrame):
    def f(rows: pd.DataFrame):
        if len(rows) != 2:
            raise ValueError("each player must have 2 stat rows for each side")
        def_stats = rows[rows["side"] == "def"].iloc[0]
        atk_stats = rows[rows["side"] == "atk"].iloc[0]

        total_kills = def_stats["kills"] + atk_stats["kills"]
        total_deaths = def_stats["deaths"] + atk_stats["deaths"]
        total_assists = def_stats["assists"] + atk_stats["assists"]
        total_rounds = len(rounds_df)
        def_rounds = len(rounds_df[rounds_df["def"] == 1])
        atk_rounds = total_rounds - def_rounds

        return pd.Series(
            {
                "agent_id": def_stats["agent_id"],
                "role": def_stats["role"],
                "def_rating": def_stats["rating"],
                "atk_rating": atk_stats["rating"],
                "rating": 0 if total_rounds == 0 else (def_stats["rating"] * def_rounds + atk_stats["rating"] * atk_rounds) / total_rounds,
                "def_acs": def_stats["acs"],
                "atk_acs": atk_stats["acs"],
                "acs": 0 if total_rounds == 0 else (def_stats["acs"] * def_rounds + atk_stats["acs"] * atk_rounds) / total_rounds,
                "def_kills_per_round": def_stats["kills"] / def_rounds,
                "atk_kills_per_round": atk_stats["kills"] / atk_rounds,
                "kills_per_round": total_kills / total_rounds,
                "def_deaths_per_round": def_stats["deaths"] / def_rounds,
                "atk_deaths_per_round": atk_stats["deaths"] / atk_rounds,
                "deaths_per_round": total_deaths / total_rounds,
                "def_assists_per_round": def_stats["assists"] / def_rounds,
                "atk_assists_per_round": atk_stats["assists"] / atk_rounds,
                "assists_per_round": total_assists / total_rounds,
                "def_kast": def_stats["kast"],
                "atk_kast": atk_stats["kast"],
                "kast": 0 if total_rounds == 0 else (def_stats["kast"] * def_rounds + atk_stats["kast"] * atk_rounds) / total_rounds,
                "def_adr": def_stats["adr"],
                "atk_adr": atk_stats["adr"],
                "adr": 0 if total_rounds == 0 else (def_stats["adr"] * def_rounds + atk_stats["adr"] * atk_rounds) / total_rounds,
                "def_hs": def_stats["hs"],
                "atk_hs": atk_stats["hs"],
                "hs": 0 if total_kills == 0 else (def_stats["hs"] * def_stats["kills"] + atk_stats["hs"] * atk_stats["kills"]) / total_kills,
                "def_fk_per_round": def_stats["first_kills"] / def_rounds,
                "atk_fk_per_round": atk_stats["first_kills"] / atk_rounds,
                "fk_per_round": (def_stats["first_kills"] + atk_stats["first_kills"]) / total_rounds,
                "def_fd_per_round": def_stats["first_deaths"] / def_rounds,
                "atk_fd_per_round": atk_stats["first_deaths"] / atk_rounds,
                "fd_per_round": (def_stats["first_deaths"] + atk_stats["first_deaths"]) / total_rounds,
                "def_1v1_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v1")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                    / def_rounds
                ),
                "atk_1v1_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v1")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "1v1_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "1v1")
                    ]
                )
                / total_rounds,
                "def_1v2_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v2")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                    / def_rounds
                ),
                "atk_1v2_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v2")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "1v2_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "1v2")
                    ]
                )
                / total_rounds,
                "def_1v3_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v3")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                    / def_rounds
                ),
                "atk_1v3_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v3")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "1v3_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "1v3")
                    ]
                )
                / total_rounds,
                "def_1v4_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v4")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                    / def_rounds
                ),
                "atk_1v4_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v4")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "1v4_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "1v4")
                    ]
                )
                / total_rounds,
                "def_1v5_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v5")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                    / def_rounds
                ),
                "atk_1v5_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "1v5")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "1v5_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "1v5")
                    ]
                )
                / total_rounds,
                "def_2ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "2k")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                )
                / def_rounds,
                "def_2ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "2k")
                        & (enriched_highlights_df["def"] == 1)
                    ]["won"].mean()
                ),
                "atk_2ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "2k")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "atk_2ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "2k")
                        & (enriched_highlights_df["def"] == 0)
                    ]["won"].mean()
                ),
                "2ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "2k")
                    ]
                )
                / total_rounds,
                "2ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "2k")
                    ]["won"].mean()
                ),
                "def_3ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "3k")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                )
                / def_rounds,
                "def_3ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "3k")
                        & (enriched_highlights_df["def"] == 1)
                    ]["won"].mean()
                ),
                "atk_3ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "3k")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "atk_3ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "3k")
                        & (enriched_highlights_df["def"] == 0)
                    ]["won"].mean()
                ),
                "3ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "3k")
                    ]
                )
                / total_rounds,
                "3ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "3k")
                    ]["won"].mean()
                ),
                "def_4ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "4k")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                )
                / def_rounds,
                "def_4ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "4k")
                        & (enriched_highlights_df["def"] == 1)
                    ]["won"].mean()
                ),
                "atk_4ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "4k")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "atk_4ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "4k")
                        & (enriched_highlights_df["def"] == 0)
                    ]["won"].mean()
                ),
                "4ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "4k")
                    ]
                )
                / total_rounds,
                "4ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "4k")
                    ]["won"].mean()
                ),
                "def_5ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "5k")
                        & (enriched_highlights_df["def"] == 1)
                    ]
                )
                / def_rounds,
                "def_5ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "5k")
                        & (enriched_highlights_df["def"] == 1)
                    ]["won"].mean()
                ),
                "atk_5ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "5k")
                        & (enriched_highlights_df["def"] == 0)
                    ]
                )
                / atk_rounds,
                "atk_5ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"])
                        & (enriched_highlights_df["highlight_type"] == "5k")
                        & (enriched_highlights_df["def"] == 0)
                    ]["won"].mean()
                ),
                "5ks_per_round": len(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "5k")
                    ]
                )
                / total_rounds,
                "5ks_convert_rate": np.nan_to_num(
                    enriched_highlights_df[
                        (enriched_highlights_df["player_id"] == def_stats["player_id"]) & (enriched_highlights_df["highlight_type"] == "5k")
                    ]["won"].mean()
                ),
            }
        )

    return f


def get_players_stats(players_stats_df: pd.DataFrame, rounds_df: pd.DataFrame, highlights_df: pd.DataFrame):
    """
    Compute the stats for each players. both players_stats_df, highlights_df and rounds_df contains data from a specific match for a specific team
    """

    try:
        enriched_highlights_df = (
            highlights_df.merge(rounds_df[["round_no", "def", "won"]], on="round_no", how="inner")
            .groupby(["match_id", "map_id", "team_id", "player_id", "round_no", "highlight_type"])
            .apply(lambda rows: rows.iloc[0])
            .reset_index(drop=True)
        )

        computed_players_stats_df = pd.DataFrame(index=range(5))
        apply_func = get_player_stats(enriched_highlights_df, rounds_df)
        computed_players_stats_df = players_stats_df.groupby("player_id").apply(apply_func)
        computed_players_stats_df.fillna(0)

        return computed_players_stats_df
    except Exception as e:
        raise e
