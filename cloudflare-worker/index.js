
export default {

    async queue(batch, env){

        for(const message of batch.message){

            const eventData = message.body;
            try{
                const response = await fetch(env.WEBHOOK_URL,{
                    method: "POST",
                    headers:{
                        "Content-Type" : "application/json",
                        "Webhook-Secret" : env.WEBHOOK_SECRET,
                    },
                    body: JSON.stringify(eventData),
                });

                if(response.ok) message.ack();
                else message.retry();

            }catch(error){
                message.retry();
            };
        };
    },
};