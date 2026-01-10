# Authentication And Load Balancer

for this service we are using rust cause its thread safe and it dosent have a guarbage collector and hepls in the long run


## What does it do ?
        it checks if the user is registered and manages the session of the user and also handles the accessibility of the use (dosent allow to query data for another user)

        this is handeled by the token provided by the auth service.

        using redis for stroing tokens os users and validating them



        redis-cli ping
