import sqlite3
import os
import numpy as np
import pandas as pd
from data_loader import load_matches
from features import avg_opps_rating_diff, direct_hth, fk_fd_per_round_diff, indirect_hth, wr_diff, maps_strength_diff

db_path = os.path.join(os.getcwd(), "../../../database/vlr.db")
dataset_path = os.path.join(os.getcwd(), "../../data/dataset.csv")

conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
conn.execute("PRAGMA journal_mode=WAL;")
conn.execute("PRAGMA synchronous=NORMAL;")
conn.execute("PRAGMA cache_size=160000;")

df = load_matches(conn)
df = df[:1000]

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


df = add_n_filter(df, "rating_diff", lambda row: (row["team_1_rating"] ** 2 - row["team_2_rating"] ** 2) / 2)
df = add_n_filter(df, "wr_diff", lambda row: wr_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]))
df = add_n_filter(df, "avg_opps_rating_diff", lambda row: avg_opps_rating_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]))
df = add_n_filter(df, "direct_hth", lambda row: direct_hth(conn, row["team_1_id"], row["team_2_id"], row["date"]))
df = add_n_filter(df, "indirect_hth", lambda row: indirect_hth(conn, row["team_1_id"], row["team_2_id"], row["date"]))
df = add_n_filter(df, "maps_strength_diff", lambda row: maps_strength_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]))
df = add_n_filter(
    df,
    ["fk_per_round_diff", "fd_per_round_diff"],
    lambda row: fk_fd_per_round_diff(conn, row["team_1_id"], row["team_2_id"], row["date"]),
)

df.to_csv(dataset_path)
print(df)
