import sqlite3
import os
import numpy as np
import pandas as pd
from data_loader import load_matches
from features import (
    avg_opps_rating_diff,
    avg_rounds_after_win_n_loss_diff,
    highlights_diff,
    direct_hth,
    fk_fd_per_round_diff,
    indirect_hth,
    key_round_wr_diff,
    round_wr_based_on_buy_type,
    thrifty_chance,
    wr_based_on_lead_diff,
    wr_diff,
    maps_strength_diff,
)

db_path = os.path.join(os.getcwd(), "../../../database/vlr.db")
dataset_path = os.path.join(os.getcwd(), "../../data/dataset_v2.csv")

conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
conn.execute("PRAGMA journal_mode=WAL;")
conn.execute("PRAGMA synchronous=NORMAL;")
conn.execute("PRAGMA cache_size=320000;")

df = load_matches(conn)

# Adding labels
df["team_won"] = np.where(df["team_1_score"] > df["team_2_score"], 0, 1)
df["result"] = df["team_1_score"].astype(str) + "-" + df["team_2_score"].astype(str)


# Adding features
def add_n_filter(df: pd.DataFrame, col_names, apply_func):
    if isinstance(col_names, str):
        df[col_names] = df.apply(apply_func, axis=1)
        df = df[df[col_names].notna()]
    elif isinstance(col_names, list):
        df[col_names] = df.apply(apply_func, axis=1, result_type="expand")
        for col_name in col_names:
            df = df[df[col_name].notna()]

    return df


df = add_n_filter(
    df,
    ["t1_wr", "t2_wr"],
    lambda row: wr_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)
df = add_n_filter(
    df,
    ["t1_avg_opps_rating", "t2_avg_opps_rating"],
    lambda row: avg_opps_rating_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    ["t1_wins_vs_t2", "t2_wins_vs_t1"],
    lambda row: direct_hth(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)
df = add_n_filter(
    df,
    ["t1_indirect_wins_vs_t2", "t2_indirect_wins_vs_t1"],
    lambda row: indirect_hth(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)
df = add_n_filter(
    df,
    "maps_strength_diff",
    lambda row: maps_strength_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    ["t1_fks_per_round", "t1_fds_per_round", "t2_fks_per_round", "t2_fds_per_round"],
    lambda row: fk_fd_per_round_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    [
        "t1_1v1s_per_round",
        "t1_1v2s_per_round",
        "t1_1v3s_per_round",
        "t1_1v4s_per_round",
        "t1_1v5s_per_round",
        "t1_2ks_per_round",
        "t1_3ks_per_round",
        "t1_4ks_per_round",
        "t1_5ks_per_round",
        "t1_2ks_converted_rate",
        "t1_3ks_converted_rate",
        "t1_4ks_converted_rate",
        "t1_5ks_converted_rate",
        "t2_1v1s_per_round",
        "t2_1v2s_per_round",
        "t2_1v3s_per_round",
        "t2_1v4s_per_round",
        "t2_1v5s_per_round",
        "t2_2ks_per_round",
        "t2_3ks_per_round",
        "t2_4ks_per_round",
        "t2_5ks_per_round",
        "t2_2ks_converted_rate",
        "t2_3ks_converted_rate",
        "t2_4ks_converted_rate",
        "t2_5ks_converted_rate",
    ],
    lambda row: highlights_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)
df = add_n_filter(
    df,
    [
        "t1_avg_rounds_win_after_round_win",
        "t1_avg_rounds_loss_after_round_loss",
        "t2_avg_rounds_win_after_round_win",
        "t2_avg_rounds_loss_after_round_loss",
    ],
    lambda row: avg_rounds_after_win_n_loss_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    ["t1_wr_with_lead", "t1_wr_without_lead", "t2_wr_with_lead", "t2_wr_without_lead"],
    lambda row: wr_based_on_lead_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    ["t1_key_round_wr", "t2_key_round_wr"],
    lambda row: key_round_wr_diff(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    [
        "t1_pistol_round_wr",
        "t1_eco_wr",
        "t1_semi_eco_wr",
        "t1_semi_buy_wr",
        "t1_full_buy_wr",
        "t2_pistol_round_wr",
        "t2_eco_wr",
        "t2_semi_eco_wr",
        "t2_semi_buy_wr",
        "t2_full_buy_wr",
    ],
    lambda row: round_wr_based_on_buy_type(
        conn, row["team_1_id"], row["team_2_id"], row["date"]
    ),
)
df = add_n_filter(
    df,
    [
        "t1_thrifty_chance",
        "t1_thriftied_chance",
        "t2_thrifty_chance",
        "t2_thriftied_chance",
    ],
    lambda row: thrifty_chance(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)

df.to_csv(dataset_path)
print(df)
