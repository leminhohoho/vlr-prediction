import math
import torch
import torch.nn as nn
import torch.nn.functional as F
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
import matplotlib.pyplot as plt

scaler = StandardScaler()

# df = pd.read_csv(os.path.join(os.path.dirname(os.path.abspath(__file__)), "../../data/dataset.csv"))
df = pd.read_csv("../../data/dataset_v2.csv")

# Drop excessive columns those are not used for training
df = df.drop(
    columns=[
        "Unnamed: 0",
        "id",
        "url",
        "date",
        "tournament_id",
        "team_1_id",
        "team_2_id",
        "team_1_score",
        "team_2_score",
    ]
)

df = df.drop(columns=["result"])
# TODO: Add 1 for result prediction later


def stage_encoding(row):
    if row["stage"] == "group_stage":
        return 0
    elif row["stage"] == "playoff":
        return 1
    elif row["stage"] == "grand_final":
        return 2


df["stage"] = df.apply(stage_encoding, axis=1)

df["rating_diff"] = df.apply(lambda row: math.log((row["team_1_rating"] + 1e-8) / (row["team_2_rating"] + 1e-8)), axis=1)
df["wr_diff"] = df.apply(lambda row: math.log((row["t1_wr"] + 1e-8) / (row["t2_wr"] + 1e-8)), axis=1)
df["avg_opps_rating_diff"] = df.apply(lambda row: math.log((row["t1_avg_opps_rating"] + 1e-8) / (row["t2_avg_opps_rating"] + 1e-8)), axis=1)
df["head_to_head"] = df.apply(lambda row: math.log((row["t1_wins_vs_t2"] + 1e-8) / (row["t2_wins_vs_t1"] + 1e-8)), axis=1)
df["indirect_head_to_head"] = df.apply(lambda row: math.log((row["t1_indirect_wins_vs_t2"] + 1e-8) / (row["t2_indirect_wins_vs_t1"] + 1e-8)), axis=1)
df["fks_per_round_diff"] = df.apply(lambda row: math.log((row["t1_fks_per_round"] + 1e-8) / (row["t2_fks_per_round"] + 1e-8)), axis=1)
df["fds_per_round_diff"] = df.apply(lambda row: math.log((row["t1_fds_per_round"] + 1e-8) / (row["t2_fds_per_round"] + 1e-8)), axis=1)
df["1v1s_per_round_diff"] = df.apply(lambda row: math.log((row["t1_1v1s_per_round"] + 1e-8) / (row["t2_1v1s_per_round"] + 1e-8)), axis=1)
df["1v2s_per_round_diff"] = df.apply(lambda row: math.log((row["t1_1v2s_per_round"] + 1e-8) / (row["t2_1v2s_per_round"] + 1e-8)), axis=1)
df["1v3s_per_round_diff"] = df.apply(lambda row: math.log((row["t1_1v3s_per_round"] + 1e-8) / (row["t2_1v3s_per_round"] + 1e-8)), axis=1)
df["1v4s_per_round_diff"] = df.apply(lambda row: math.log((row["t1_1v4s_per_round"] + 1e-8) / (row["t2_1v4s_per_round"] + 1e-8)), axis=1)
df["1v5s_per_round_diff"] = df.apply(lambda row: math.log((row["t1_1v5s_per_round"] + 1e-8) / (row["t2_1v5s_per_round"] + 1e-8)), axis=1)
df["2ks_per_round_diff"] = df.apply(lambda row: math.log((row["t1_2ks_per_round"] + 1e-8) / (row["t2_2ks_per_round"] + 1e-8)), axis=1)
df["3ks_per_round_diff"] = df.apply(lambda row: math.log((row["t1_3ks_per_round"] + 1e-8) / (row["t2_3ks_per_round"] + 1e-8)), axis=1)
df["4ks_per_round_diff"] = df.apply(lambda row: math.log((row["t1_4ks_per_round"] + 1e-8) / (row["t2_4ks_per_round"] + 1e-8)), axis=1)
df["5ks_per_round_diff"] = df.apply(lambda row: math.log((row["t1_5ks_per_round"] + 1e-8) / (row["t2_5ks_per_round"] + 1e-8)), axis=1)
df["2ks_converted_rate_diff"] = df.apply(lambda row: math.log((row["t1_2ks_converted_rate"] + 1e-8) / (row["t2_2ks_converted_rate"] + 1e-8)), axis=1)
df["3ks_converted_rate_diff"] = df.apply(lambda row: math.log((row["t1_3ks_converted_rate"] + 1e-8) / (row["t2_3ks_converted_rate"] + 1e-8)), axis=1)
df["4ks_converted_rate_diff"] = df.apply(lambda row: math.log((row["t1_4ks_converted_rate"] + 1e-8) / (row["t2_4ks_converted_rate"] + 1e-8)), axis=1)
df["5ks_converted_rate_diff"] = df.apply(lambda row: math.log((row["t1_5ks_converted_rate"] + 1e-8) / (row["t2_5ks_converted_rate"] + 1e-8)), axis=1)
df["avg_rounds_win_after_round_win_diff"] = df.apply(
    lambda row: math.log((row["t1_avg_rounds_win_after_round_win"] + 1e-8) / (row["t2_avg_rounds_win_after_round_win"] + 1e-8)), axis=1
)
df["avg_rounds_loss_after_round_loss_diff"] = df.apply(
    lambda row: -1 * math.log((row["t1_avg_rounds_loss_after_round_loss"] + 1e-8) / (row["t2_avg_rounds_loss_after_round_loss"] + 1e-8)), axis=1
)
df["wr_with_lead_diff"] = df.apply(lambda row: math.log((row["t1_wr_with_lead"] + 1e-8) / (row["t2_wr_with_lead"] + 1e-8)), axis=1)
df["wr_without_lead_diff"] = df.apply(lambda row: math.log((row["t1_wr_without_lead"] + 1e-8) / (row["t2_wr_without_lead"] + 1e-8)), axis=1)
df["key_round_wr_diff"] = df.apply(lambda row: math.log((row["t1_key_round_wr"] + 1e-8) / (row["t2_key_round_wr"] + 1e-8)), axis=1)
df["pistol_round_wr_diff"] = df.apply(lambda row: math.log((row["t1_pistol_round_wr"] + 1e-8) / (row["t2_pistol_round_wr"] + 1e-8)), axis=1)
df["eco_wr_diff"] = df.apply(lambda row: math.log((row["t1_eco_wr"] + 1e-8) / (row["t2_eco_wr"] + 1e-8)), axis=1)
df["semi_eco_wr_diff"] = df.apply(lambda row: math.log((row["t1_semi_eco_wr"] + 1e-8) / (row["t2_semi_eco_wr"] + 1e-8)), axis=1)
df["semi_buy_wr_diff"] = df.apply(lambda row: math.log((row["t1_semi_buy_wr"] + 1e-8) / (row["t2_semi_buy_wr"] + 1e-8)), axis=1)
df["full_buy_wr_diff"] = df.apply(lambda row: math.log((row["t1_full_buy_wr"] + 1e-8) / (row["t2_full_buy_wr"] + 1e-8)), axis=1)

df = df.drop(
    columns=[
        # "4ks_per_round_diff",
        # "t1_4ks_per_round",
        # "t2_4ks_per_round",
        # "5ks_per_round_diff",
        # "t1_5ks_per_round",
        # "t2_5ks_per_round",
        # "3ks_converted_rate_diff",
        # "t1_3ks_converted_rate",
        # "t2_3ks_converted_rate",
        # "5ks_converted_rate_diff",
        # "t1_5ks_converted_rate",
        # "t2_5ks_converted_rate",
        # "1v1s_per_round_diff",
        # "t1_1v1s_per_round",
        # "t2_1v1s_per_round",
        # "1v5s_per_round_diff",
        # "t1_1v5s_per_round",
        # "t2_1v5s_per_round",
        # "head_to_head",
        # "t1_wins_vs_t2",
        # "t2_wins_vs_t1",
        # "indirect_head_to_head",
        # "t1_indirect_wins_vs_t2",
        # "t2_indirect_wins_vs_t1",
    ]
)


# Model for predicting team won
class TeamWonModel(nn.Module):
    def __init__(self, in_features=97, hidden=64, out_features=2):
        super().__init__()
        self.network = nn.Sequential(
            nn.Linear(in_features, hidden),
            # nn.Linear(hidden, hidden),
            # nn.Linear(hidden, hidden),
            # nn.Linear(hidden, hidden),
            # nn.Linear(hidden, hidden),
            nn.Linear(hidden, out_features),
        )

    def forward(self, x):
        return self.network(x)


model = TeamWonModel()
torch.manual_seed(42)

# Split data into train and test
X = df.drop("team_won", axis=1)
X = scaler.fit_transform(X)
y = df["team_won"]

y = y.values

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)
X_train = torch.FloatTensor(X_train)
X_test = torch.FloatTensor(X_test)
y_train = torch.LongTensor(y_train)
y_test = torch.LongTensor(y_test)


criterion = nn.CrossEntropyLoss()
# criterion = nn.BCEWithLogitsLoss()
optimizer = torch.optim.Adam(model.parameters(), lr=1e-3)

# Training the model
epochs = 10000  # NOTE: Training an actual model for production will need more epochs, this is just for testing
losses = []
for i in range(epochs):
    y_pred = model.forward(X_train)

    loss = criterion(y_pred, y_train)
    losses.append(loss.detach().numpy())

    if (i + 1) % 10 == 0:
        print(f"Epoch {i+1} and loss: {loss}")

    optimizer.zero_grad()
    loss.backward()
    optimizer.step()

with torch.no_grad():
    y_eval = model.forward(X_test)
    loss = criterion(y_eval, y_test)

    print(f"loss: {loss}")

    corrects = 0

    for i, data in enumerate(X_test):
        y_val = model.forward(data)

        if y_val.argmax().item() == y_test[i]:
            corrects += 1

    print(f"Corrected: {corrects}")
    print(f"Wrong: {X_test.__len__() - corrects}")
    print(f"Accuracy: {corrects/len(X_test)}")


plt.plot(range(epochs), losses)
plt.ylabel("loss/error")
plt.xlabel("epochs")
plt.savefig("loss_graph.png")
plt.show()
