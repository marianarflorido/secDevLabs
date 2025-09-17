var db = connect("mongodb://localhost/DB");
db.createUser(
    {
        user: "User2806622879",
        pwd: "Pass790923903",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);
