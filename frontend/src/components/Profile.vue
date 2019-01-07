<template>
  <div class="about">
    <h1>Welcome {{ username }}</h1>
    <p>{{ email }}</p>
  </div>
</template>
<script>
export default {
    name: 'about',
    data() {
        return {
            username: '',
            email: '',
        };
    },
    props: {
        authenticated: Boolean,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.getProtected('/api/user')
            .then((response) => {
                this.username = response.data.username;
                this.email = response.data.email;
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
