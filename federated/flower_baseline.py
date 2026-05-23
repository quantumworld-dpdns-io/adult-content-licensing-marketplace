"""Flower baseline placeholder for federated recommendation training."""

def train_round(metrics):
    return {"loss": metrics.get("loss", 1.0) * 0.95}
