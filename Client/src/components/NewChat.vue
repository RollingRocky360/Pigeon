<script>
    export default {
        inject: ['user'],
        data() {
            return {
                searchTerm: '',
                searchUser: undefined
            }
        },
        methods: {
            async search() {
                if (this.searchTerm.length === 0) return;
                const resp = await fetch('http://localhost:3000/user/' + this.searchTerm);
                this.searchUser = await resp.json()
            },
            async testEnterAndSearch(e) {
                if (e.key !== 'Enter') return;
                this.search()
            },
            async createChat(otherUser) {
                const resp = await fetch('http://localhost:3000/chat', {
                    method: 'POST',
                    body: JSON.stringify([this.user.id, otherUser.id])
                })
                if (resp.ok) {
                    const chat = await resp.json()
                    this.$emit('newchat', chat)
                }
            }
        }
    }
</script>

<template>
    <div id="search-area">
        <div id="search-input-area">
            <input type="text" v-model="searchTerm" id="search" @keyup="testEnterAndSearch" placeholder="Search for a friend...">
            <button @click="search" title="search">
                <span class="material-symbols-outlined">
                    search
                </span>
            </button>
        </div>
        <div v-if="searchUser === null">
            No User Found
        </div>
        <div v-else-if="searchUser !== undefined" 
            id="search-user-details" @click="() => createChat(searchUser)"
            title="Add">
            <div id="pfp">{{ searchUser.username[0] }}</div>
            <p id="username">{{ searchUser.username }}</p>
            <p id="email">{{ searchUser.email }}</p>
        </div>
    </div>
</template>

<style scoped>
    #search-area {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    #search-user-details {
        display: grid;
        grid-template-columns: auto auto;
        grid-template-rows: auto auto;
        place-items: center;
        column-gap: .5rem;
        margin: 1rem 0;
        cursor: pointer;
        padding: 1rem 1.5rem;
        border-radius: 1rem;
    }
    #search-user-details:hover {
        background-color: #00000010;
    }

    #search-input-area {
        display: flex;
        width: 60%;
        justify-content: space-between;
        align-items: center;
        border: 2px solid #72737535;
        padding: .5rem .7rem .5rem 1.4rem;
        border-radius: 5rem;
        background-color: white;
    }

    #pfp {
        width: 3.5rem;
        font-size: 1.1rem;
        aspect-ratio: 1;
        grid-row: 1/3;
        border-radius: 50%;
        border: 1px solid #00000040;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: white;
    }
    #username, #email {
        justify-self: left;
    }
    #username {
        font-size: 1.1rem;
    }
    #email {
        color: #00000060;
        font-size: .9rem;
        margin-top: -.7rem;
    }

    input {
        all: unset;
        width: 90%;
    }
    input::placeholder {
        color: #00000050;
    }
    input:focus::placeholder {
        color: transparent;
    }
    button {
        all: unset;
        padding: .2rem;
        display: flex;
        align-items: center;
        cursor: pointer;
    }
    button:hover {
        color: #00000060;
    }

</style>