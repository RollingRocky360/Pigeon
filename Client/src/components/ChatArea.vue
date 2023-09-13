<script>
    export default {
        inject: [ 'websocket', 'user', 'currentChat' ],
        props: ['chatId', 'currentChat'],
        data() {
            return {
                messages: [],
                newMessage: "",
                chat: null,
            }
        },
        methods: {
            checkEnterAndSend(e) {
                if (e.key !== 'Enter') return
                this.sendMessage()
            },
            sendMessage() {
                if (this.newMessage.length === 0) return
                this.websocket.send(JSON.stringify({
                    type: 'new',
                    message: {
                        "chat_id": this.chatId, 
                        "user_id": this.user.id,
                        "content": this.newMessage
                    }
                }))
                this.newMessage = ''
            },
            async deleteMessage(message) {
                this.websocket.send(JSON.stringify({
                    type:'delete',
                    message: message
                }))
            }
        },
        async mounted() {
            this.websocket.addEventListener('message', e => {
                if (e.data === 'error') {
                    console.error('NewMessage/DeleteMessage error')
                    return
                }
                const wsm = JSON.parse(e.data)
                if (wsm.type === 'new') this.messages.push(wsm.message)
                else if (wsm.type === 'delete') this.messages = this.messages.filter(m => m.id !== wsm.message.id)
            })
            const [respMsgs, respChat] = await Promise.all([
                fetch(`http://localhost:3000/messages/${this.chatId}`),
                fetch(`http://localhost:3000/chat/${this.chatId}`),
            ])
            if (respMsgs.ok) this.messages = await respMsgs.json()
            if (respChat.ok) this.chat = await respChat.json()
            
            console.log('user', this.user)
            console.log(this.messages)
        },
        computed: {
            messagesRev() {
                return this.messages.slice().reverse()
            }
        }
    }
</script>

<template>
    <div id="chatarea-container">
        <div id="chat-details" v-if="chat">
            <div id="pfp" :style="{ backgroundColor: chat.pfp_color }">
                {{ (chat.users[0].id !== user.id ? chat.users[0].username : chat.users[1].username)[0] }}
            </div>
            <p id="chatname">
                {{ chat.users[0].id !== user.id ? chat.users[0].username : chat.users[1].username }}
            </p>
        </div>

        <ul id="chatarea" v-if="chat">
            <li v-for="message in messagesRev" :key="message.id" :class="{ 'outgoing': message.user_id === user.id }">
                <p class="message">
                    {{ message.content }}
                </p>
                <span 
                    class="material-symbols-outlined delete" 
                    style="font-size: 1.1rem"
                    v-if="message.user_id === user.id"
                    title="delete"
                    @click="() => deleteMessage(message)">
                    delete
                </span>
                <span class="datetime">
                    {{ new Date(parseInt(message.timestamp)/1000000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
                </span>
            </li>
        </ul>

        <div id="actions" v-if="chat">
            <input type="text" v-model="newMessage" placeholder="Type a message" id="newmessage" @keyup="checkEnterAndSend">
            <button @click="sendMessage" id="send">
                <span class="material-symbols-outlined" style="font-size: 1rem;">
                    send
                </span>
            </button>
        </div>
    </div>  
</template>

<style scoped>
    #chatarea-container {
        width: 99%;
        height: 95%;
        background-color: #ffffff;
        border-radius: .5rem;
        box-sizing: border-box;
        padding: 1rem;
        display: grid;
        gap: .5rem;
        grid-template-rows: 3rem .9fr .1fr;
        align-items: center;
        overflow: hidden;
    }
    #chatarea {
        list-style-type: none;    
        padding: 0 1rem;
        height: 95%;
        width: 100%;
        display: flex;
        gap: .25rem;
        flex-direction: column-reverse;
        overflow-y: auto;
    }

    #newmessage {
        all: unset;
        height: 2rem;
        width: 90%;
        font-size: 1rem;
        padding: 0 .2rem;
        box-sizing: border-box;
        border-radius: 5rem;
    }
    #newmessage::placeholder {
        color: #7273756c;
    }
    #newmessage:focus::placeholder {
        color: transparent;
    }

    #actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        border: 1px solid #72737535;
        padding: .5rem .7rem .5rem 1rem;
        border-radius: 5rem;
    }

    #send {
        all: unset;
        background-color: #9090ff;
        color: white;
        width: 2rem;
        height: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        transition-duration: .2s;
        cursor: pointer;
    }

    #chat-details {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        box-sizing: border-box;
        padding: .5rem;
        gap: .7rem;
        border-bottom: 1px solid #00000020;
    }

    #pfp {
        font-size: 1.5rem;
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        height: 2.5rem;
        border-radius: 50%;
    }

    #chatname {
        font-size: 1.15rem;
    }


    li > .delete, .datetime {
        color: transparent;
        transition: .2s;
    }
    li:hover > .delete {
        color: #00000070;
    }
    li:hover > .datetime {
        color: #00000070;
    }

    li {
        max-width: 75%;
        height: max-content;
        width: max-content;
        cursor: default;
        display: flex;
        align-items: center;
        gap: .7rem;
    }

    li > .message {
        padding: .3rem 1rem;
        border: 1px solid #00000020;
        border-radius: 1rem;
    }

    li > .datetime {
        text-wrap: nowrap;
        font-size: .7rem;
    }
    .outgoing {
        align-self: flex-end;
        flex-direction: row-reverse;
    }
    .outgoing > .message {
        background-color: #9090ff;
        color: white;
        border: none;
    }

    .delete {
        aspect-ratio: 1;
        padding: .3rem;
        border-radius: 50%;
        user-select: none;
    }
    .delete:hover {
        background-color: #00000020;
    }

</style>