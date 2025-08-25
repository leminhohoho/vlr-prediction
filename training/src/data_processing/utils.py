from datetime import datetime


def subtract_date(date_str_1, date_str_2):
    datetime_1 = datetime.strptime(date_str_1, "%Y-%m-%d %H:%M:%S%z")
    datetime_2 = datetime.strptime(date_str_2, "%Y-%m-%d %H:%M:%S%z")

    return (datetime_1 - datetime_2).days


def calc_wr(wins, losses):
    if wins == 0 and losses == 0:
        return 0

    return wins / (wins + losses)


def divide(a, b):
    return 0 if b == 0 else a / b
