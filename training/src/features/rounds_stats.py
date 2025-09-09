import pandas as pd


def compute_round_stats(clutches_df: pd.DataFrame, opp_clutches_df: pd.DataFrame):
    def f(row: pd.Series):
        return pd.Series(
            {
                "round_no": row["round_no"],
                "def": row["def"],
                "team_buy_type": row["team_buy_type"],
                "team_against_buy_type": row["team_against_buy_type"],
                "team_bank": row["team_bank"],
                "team_against_bank": row["team_against_bank"],
                "won": row["won"],
                "won_buy_elim": row["won"] == 1 and row["won_method"] == "eliminate",
                "loss_buy_elim": row["won"] == 0 and row["won_method"] == "eliminate",
                "won_buy_spike_explode": row["won"] == 1 and row["won_method"] == "spike_explode",
                "loss_buy_spike_explode": row["won"] == 0 and row["won_method"] == "spike_explode",
                "won_buy_defuse": row["won"] == 1 and row["won_method"] == "defuse",
                "loss_buy_defuse": row["won"] == 0 and row["won_method"] == "defuse",
                "won_buy_timeout": row["won"] == 1 and row["won_method"] == "out_of_time",
                "loss_buy_timeout": row["won"] == 0 and row["won_method"] == "out_of_time",
                "won_by_thrifty": row["won"] == 1
                and (
                    row["team_buy_type"] == "eco"
                    and row["team_against_buy_type"] in ["semi_buy", "full_buy"]
                    or (row["team_buy_type"] == "semi_eco" and row["team_against_buy_type"] == "full_buy")
                ),
                "loss_by_thrifty": row["won"] == 0
                and (
                    row["team_against_buy_type"] == "eco"
                    and row["team_buy_type"] in ["semi_buy", "full_buy"]
                    or (row["team_against_buy_type"] == "semi_eco" and row["team_buy_type"] == "full_buy")
                ),
                "won_buy_clutch": row["won"] == 1 and not clutches_df[clutches_df["round_no"] == row["round_no"]].empty,
                "loss_buy_clutch": row["won"] == 0 and not opp_clutches_df[opp_clutches_df["round_no"] == row["round_no"]].empty,
            }
        )

    return f


def get_rounds_stats(rounds_df: pd.DataFrame, highlights_df: pd.DataFrame, opp_highlights_df: pd.DataFrame):
    def get_clutches(highlights_df: pd.DataFrame):
        return (
            highlights_df[highlights_df["highlight_type"].isin(["1v1", "1v2", "1v3", "1v4", "1v5"])]
            .merge(rounds_df[["round_no", "def", "won"]], on="round_no", how="inner")
            .groupby(["match_id", "map_id", "team_id", "player_id", "round_no", "highlight_type"])
            .apply(lambda rows: rows.iloc[0])
            .reset_index(drop=True)
        )

    clutches_df = get_clutches(highlights_df)
    opp_clutches_df = get_clutches(opp_highlights_df)

    return rounds_df.apply(compute_round_stats(clutches_df, opp_clutches_df), axis=1)
