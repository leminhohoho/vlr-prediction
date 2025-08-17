import math
import pandas as pd
from data_loader import (
    load_current_map_pool,
    load_team_highlights,
    load_team_fkfd,
    load_team_maps_stats_recently,
    load_team_played_maps_recently,
    load_team_played_rounds_recently,
)
from utils import subtract_date, divide


def wr_diff(conn, t1_id, t2_id, date, min_maps=16):
    t1_maps = load_team_played_maps_recently(conn, t1_id, date)
    t2_maps = load_team_played_maps_recently(conn, t2_id, date)

    if len(t1_maps) < min_maps or len(t2_maps) < min_maps:
        return None

    def append_wins(maps, team_id):
        team_wins = 0

        for map in maps.itertuples(index=False):
            date_diff = 1 - (subtract_date(date, map.date) / 1000)
            if (map.team_1_id == team_id and map.team_1_score > map.team_2_score) or (
                map.team_2_id == team_id and map.team_1_score < map.team_2_score
            ):
                team_wins += 1 * date_diff

        return team_wins

    t1_wins = append_wins(t1_maps, t1_id)
    t2_wins = append_wins(t2_maps, t2_id)

    t1_wr = t1_wins / len(t1_maps)
    t2_wr = t2_wins / len(t2_maps)

    return (t1_wr**2 - t2_wr**2) / 2


def avg_opps_rating_diff(conn, t1_id, t2_id, date, min_maps=16):
    t1_maps = load_team_played_maps_recently(conn, t1_id, date)
    t2_maps = load_team_played_maps_recently(conn, t2_id, date)

    if len(t1_maps) < min_maps or len(t2_maps) < min_maps:
        return None

    def calc_avg_opps_rating(maps, team_id):
        opps_rating = 0

        for map in maps.itertuples(index=False):
            if map.team_1_id == team_id:
                opps_rating += map.team_2_rating
            elif map.team_2_id == team_id:
                opps_rating += map.team_1_rating

        return opps_rating / len(maps)

    t1_avg_opps_rating = calc_avg_opps_rating(t1_maps, t1_id)
    t2_avg_opps_rating = calc_avg_opps_rating(t2_maps, t2_id)

    return (t1_avg_opps_rating**2 - t2_avg_opps_rating**2) / 2


def direct_hth(conn, t1_id, t2_id, date, min_maps=16):
    t1_maps = load_team_played_maps_recently(conn, t1_id, date)

    if len(t1_maps) < min_maps:
        return None

    t1_wins = 0
    t2_wins = 0

    for map in t1_maps.itertuples(index=False):
        date_diff = 1 - (subtract_date(date, map.date) / 1000)
        if map.team_1_id == t1_id and map.team_2_id == t2_id:
            if map.team_1_score > map.team_2_score:
                t1_wins += 1 * date_diff**2
            elif map.team_1_score < map.team_2_score:
                t2_wins += 1 * date_diff**2
        if map.team_1_id == t2_id and map.team_2_id == t1_id:
            if map.team_1_score > map.team_2_score:
                t2_wins += 1 * date_diff**2
            elif map.team_1_score < map.team_2_score:
                t1_wins += 1 * date_diff**2

    return (t1_wins**2 - t2_wins**2) / 2


def indirect_hth(conn, t1_id, t2_id, date):
    t1_maps = load_team_played_maps_recently(conn, t1_id, date)
    t2_maps = load_team_played_maps_recently(conn, t2_id, date)

    if len(t1_maps) < 16 or len(t2_maps) < 16:
        return None

    t1_wins = 0
    t2_wins = 0

    teams_compared = []

    for map in t1_maps.itertuples(index=False):
        if map.team_1_id == t1_id:
            t_intermediate = map.team_2_id
        elif map.team_2_id == t1_id:
            t_intermediate = map.team_1_id

        if t_intermediate in teams_compared or t2_maps[(t2_maps["team_1_id"] == t_intermediate) | (t2_maps["team_2_id"] == t_intermediate)].empty:
            continue

        t1_vs_t_intermediate = direct_hth(conn, t1_id, t_intermediate, date)
        t2_vs_t_intermediate = direct_hth(conn, t2_id, t_intermediate, date)

        if t1_vs_t_intermediate * t2_vs_t_intermediate >= 0:
            continue

        t1_wins += t1_vs_t_intermediate
        t2_wins += t2_vs_t_intermediate

    return (t1_wins**2 - t2_wins**2) / 2


def maps_strength_diff(conn, t1_id, t2_id, date):
    t1_maps_stats = load_team_maps_stats_recently(conn, t1_id, date)
    t2_maps_stats = load_team_maps_stats_recently(conn, t2_id, date)
    current_maps_pool = load_current_map_pool(conn, date)

    strength_diff = 0

    for map in current_maps_pool.itertuples(index=False):
        t1_map_stats = t1_maps_stats[t1_maps_stats.map_id == map.map_id]
        t2_map_stats = t2_maps_stats[t2_maps_stats.map_id == map.map_id]

        t1_map_wr = 0 if t1_map_stats.empty else divide(t1_map_stats.iloc[0]["wins"], t1_map_stats.iloc[0]["losses"])
        t2_map_wr = 0 if t2_map_stats.empty else divide(t2_map_stats.iloc[0]["wins"], t2_map_stats.iloc[0]["losses"])

        map_wr_diff = (t1_map_wr**2 - t2_map_wr**2) / 2
        map_wr_diff *= (1 - math.sqrt(map_wr_diff**2)) ** 2

        strength_diff += map_wr_diff

    return strength_diff


def fk_fd_per_round_diff(conn, t1_id, t2_id, date, min_maps=16):
    t1_fkfds = load_team_fkfd(conn, t1_id, date)
    t2_fkfds = load_team_fkfd(conn, t2_id, date)

    if len(t1_fkfds) < min_maps or len(t2_fkfds) < min_maps:
        return (None, None)

    t1_fk_per_rounds = t1_fkfds["fks"].sum() / t1_fkfds["rounds"].sum()
    t1_fd_per_rounds = t1_fkfds["fds"].sum() / t1_fkfds["rounds"].sum()
    t2_fk_per_rounds = t2_fkfds["fks"].sum() / t2_fkfds["rounds"].sum()
    t2_fd_per_rounds = t2_fkfds["fds"].sum() / t2_fkfds["rounds"].sum()

    return ((t1_fk_per_rounds**2 - t2_fk_per_rounds**2) / 2, (t1_fd_per_rounds**2 - t2_fd_per_rounds**2) / 2)


def highlights_diff(conn, t1_id, t2_id, date, min_maps=16):
    try:
        t1_highlights_log = load_team_highlights(conn, t1_id, date)
        t2_highlights_log = load_team_highlights(conn, t2_id, date)
        t1_maps = load_team_played_maps_recently(conn, t1_id, date)
        t2_maps = load_team_played_maps_recently(conn, t2_id, date)
        t1_rounds = load_team_played_rounds_recently(conn, t1_id, date)
        t2_rounds = load_team_played_rounds_recently(conn, t2_id, date)

        t1_rounds_played = t1_maps["team_1_score"].sum() + t1_maps["team_2_score"].sum()
        t2_rounds_played = t2_maps["team_1_score"].sum() + t2_maps["team_2_score"].sum()

        if len(t1_maps) < min_maps or len(t2_maps) < min_maps or t1_rounds_played != len(t1_rounds) or t2_rounds_played != len(t2_rounds):
            raise Exception()

        def is_mk_win(rounds: pd.DataFrame, hl_type):
            def convered(row):
                round = rounds[
                    (rounds["round_no"] == row["round_no"]) & (rounds["match_id"] == row["match_id"]) & (rounds["map_id"] == row["map_id"])
                ]
                if round.empty:
                    raise Exception()
                return row["team_id"] == round.iloc[0]["team_won"] and row["highlight_type"] == hl_type

            return convered

        t1_1v1s = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "1v1"])
        t1_1v2s = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "1v2"]) / 2
        t1_1v3s = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "1v3"]) / 3
        t1_1v4s = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "1v4"]) / 4
        t1_1v5s = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "1v5"]) / 5
        t1_2ks = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "2k"]) / 2
        t1_3ks = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "3k"]) / 3
        t1_4ks = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "4k"]) / 4
        t1_5ks = len(t1_highlights_log[t1_highlights_log["highlight_type"] == "5k"]) / 5
        t1_2ks_converted = divide(len(t1_highlights_log[t1_highlights_log.apply(is_mk_win(t1_rounds, "2k"), axis=1)]) / 2, t1_2ks)
        t1_3ks_converted = divide(len(t1_highlights_log[t1_highlights_log.apply(is_mk_win(t1_rounds, "3k"), axis=1)]) / 3, t1_2ks)
        t1_4ks_converted = divide(len(t1_highlights_log[t1_highlights_log.apply(is_mk_win(t1_rounds, "4k"), axis=1)]) / 4, t1_2ks)
        t1_5ks_converted = divide(len(t1_highlights_log[t1_highlights_log.apply(is_mk_win(t1_rounds, "5k"), axis=1)]) / 5, t1_2ks)
        t2_1v1s = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "1v1"])
        t2_1v2s = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "1v2"]) / 2
        t2_1v3s = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "1v3"]) / 3
        t2_1v4s = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "1v4"]) / 4
        t2_1v5s = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "1v5"]) / 5
        t2_2ks = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "2k"]) / 2
        t2_3ks = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "3k"]) / 3
        t2_4ks = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "4k"]) / 4
        t2_5ks = len(t2_highlights_log[t2_highlights_log["highlight_type"] == "5k"]) / 5
        t2_2ks_converted = divide(len(t2_highlights_log[t2_highlights_log.apply(is_mk_win(t2_rounds, "2k"), axis=1)]) / 2, t2_2ks)
        t2_3ks_converted = divide(len(t2_highlights_log[t2_highlights_log.apply(is_mk_win(t2_rounds, "3k"), axis=1)]) / 3, t2_2ks)
        t2_4ks_converted = divide(len(t2_highlights_log[t2_highlights_log.apply(is_mk_win(t2_rounds, "4k"), axis=1)]) / 4, t2_2ks)
        t2_5ks_converted = divide(len(t2_highlights_log[t2_highlights_log.apply(is_mk_win(t2_rounds, "5k"), axis=1)]) / 5, t2_2ks)

        return (
            ((t1_1v1s / t1_rounds_played) ** 2 - (t2_1v1s / t2_rounds_played) ** 2) / 2,
            ((t1_1v2s / t1_rounds_played) ** 2 - (t2_1v2s / t2_rounds_played) ** 2) / 2,
            ((t1_1v3s / t1_rounds_played) ** 2 - (t2_1v3s / t2_rounds_played) ** 2) / 2,
            ((t1_1v4s / t1_rounds_played) ** 2 - (t2_1v4s / t2_rounds_played) ** 2) / 2,
            ((t1_1v5s / t1_rounds_played) ** 2 - (t2_1v5s / t2_rounds_played) ** 2) / 2,
            ((t1_2ks / t1_rounds_played) ** 2 - (t2_2ks / t2_rounds_played) ** 2) / 2,
            ((t1_3ks / t1_rounds_played) ** 2 - (t2_3ks / t2_rounds_played) ** 2) / 2,
            ((t1_4ks / t1_rounds_played) ** 2 - (t2_4ks / t2_rounds_played) ** 2) / 2,
            ((t1_5ks / t1_rounds_played) ** 2 - (t2_5ks / t2_rounds_played) ** 2) / 2,
            (t1_2ks_converted**2 - t2_2ks_converted**2) / 2,
            (t1_3ks_converted**2 - t2_3ks_converted**2) / 2,
            (t1_4ks_converted**2 - t2_4ks_converted**2) / 2,
            (t1_5ks_converted**2 - t2_5ks_converted**2) / 2,
        )
    except Exception as e:
        return (None, None, None, None, None, None, None, None, None, None, None, None, None)


def avg_rounds_after_win_n_loss(conn, t1_id, t2_id, date, min_maps=16):
    t1_rounds = load_team_played_rounds_recently(conn, t1_id, date)
    t2_rounds = load_team_played_rounds_recently(conn, t2_id, date)
    t1_maps = load_team_played_maps_recently(conn, t1_id, date)
    t2_maps = load_team_played_maps_recently(conn, t2_id, date)

    t1_rounds_played = t1_maps["team_1_score"].sum() + t1_maps["team_2_score"].sum()
    t2_rounds_played = t2_maps["team_1_score"].sum() + t2_maps["team_2_score"].sum()

    if len(t1_maps) < min_maps or len(t2_maps) < min_maps or len(t1_rounds) != t1_rounds_played or len(t2_rounds) != t2_rounds_played:
        return (None, None)

    def avg(rounds: pd.DataFrame, team_id):
        maps = rounds.groupby(["match_id", "map_id"])
        rounds_win = win_streaks = rounds_loss = loss_streaks = 0
        for _, m in maps:
            rwc = rlc = 0
            for r in m.itertuples(index=False):
                if r.team_won == team_id:
                    rwc += 1
                    if rlc:
                        rounds_loss += rlc - 1
                        loss_streaks += 1
                        rlc = 0
                else:
                    rlc += 1
                    if rwc:
                        rounds_win += rwc - 1
                        win_streaks += 1
                        rwc = 0
            if rwc:
                rounds_win += rwc - 1
                win_streaks += 1
            if rlc:
                rounds_loss += rlc - 1
                loss_streaks += 1

        return (divide(rounds_win, win_streaks), divide(rounds_loss, loss_streaks))

    t1_avg_round_win_after_win, t1_avg_round_loss_after_loss = avg(t1_rounds, t1_id)
    t2_avg_round_win_after_win, t2_avg_round_loss_after_loss = avg(t2_rounds, t2_id)

    return (
        (t1_avg_round_win_after_win**2 - t2_avg_round_win_after_win**2) / 2,
        (t1_avg_round_loss_after_loss**2 - t2_avg_round_loss_after_loss**2) / 2,
    )
