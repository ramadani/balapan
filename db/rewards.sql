CREATE TABLE "public"."rewards" (
    "id" uuid NOT NULL,
    "transaction_limit" numeric NOT NULL,
    "transaction_usage" numeric NOT NULL,
    "rewards_limit" numeric NOT NULL,
    "rewards_usage" numeric NOT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."rewards" ("id", "transaction_limit", "transaction_usage", "rewards_limit", "rewards_usage") VALUES
('51cc102d-4583-48f3-b8b1-2073b9c0663f', 75, 0, 100000, 0),
('8470581f-18ba-465b-9933-d44e12679b45', 100, 0, 75000, 0),
('d5b57ea6-468e-4f73-9bd2-d062cfeb5109', 100, 0, 100000, 0);