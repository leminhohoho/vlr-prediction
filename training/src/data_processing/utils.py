from datetime import datetime


def subtract_date(date_str_1, date_str_2):
    datetime_1 = datetime.strptime(date_str_1, "%Y-%m-%d %H:%M:%S%z")
    datetime_2 = datetime.strptime(date_str_2, "%Y-%m-%d %H:%M:%S%z")

    return (datetime_1 - datetime_2).days
