# IsaacRules Service

Este serviço tem como objetivo principal disponibilizar crédito e antecipação(Operation Domain) para usários. Porém alem desse objetivo principal temos dominios auxiliares como o blocklist que servirá para realizar o bloqueio e desbloqueio de usuários, dominio de regras(rules) que servirá para usuários administrativos consigam realizar configurações de regras similar /Users/idsc/Projects/isaacdsc/rulesengine/jdm_graph.json.


## Blocklist Domain  
>> POST /backoffice/block-list
>> body: { "user_id": UUID }
>> save in folder ./blocklist
>> response 201

>> DELETE /backoffice/block-list/{user_id}
>> body: { "user_id": UUID }
>> remove data the load folder ./blocklist
>> response 200


## Rules Domain
>> PATCH /backoffice/rules/{campaingName}?operations=credit,anticipation
>> body: map[string]any
>> save in folder ./rules/{campaingName}

>> GET /backoffice/rules/list
>> retrieve all campins in /rules/*.json
>> output: [{"campaingx": {...}}]

>> DELETE /backoffice/{operation}/{campaingName}


## Operation Domain
>> HTTP server native golang 
>> operation: [credit,antecipation]
>> POST /rule-engine/{operation}
>> body: { "requested_amout": int(cents), "user_id": UUID } 
>> Response sucess: { "offers": [...] }
>> Bad request 400 if not existent operation, invalid requested_amout or invalid user_id
>> Unprocess entity 422 when not found user
>> Response error with status 200 when retrieve message with denied requested credit or anticipation

### Flow execution
    1. Validation parameters and body
    2. Validate if user is blocked
    3. Execute all rules in parallel
        This router read all rules /rules/*.json and merge in map[string]
        Execute all rules using goroutine and return map[campaingName]output

