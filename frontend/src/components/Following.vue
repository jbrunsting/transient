<template>
  <div class="home">
    <h1>Welcome {{ username }}</h1>
  </div>
</template>

<script>
export default {
    name: 'home',
    data() {
        return {
            id: '',
            username: '',
            email: '',
        };
    },
    props: { authenticated: Boolean },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.getProtected('/api/user')
            .then((response) => {
                this.id = response.data.id;
                this.username = response.data.username;
                this.email = response.data.email;
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
