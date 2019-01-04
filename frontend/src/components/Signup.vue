<template>
  <div class="signup">
    <Error direction="right" class="uniqueness error">
      That username or email is already taken
    </Error>
    <Error direction="right" class="format error">
      Username or email format incorrect
    </Error>
    <Error direction="right" class="unknown error">
      Could not sign up, please try again later
    </Error>
    <form @submit.prevent="signup">
      <input placeholder="username" v-model="username">
      <input type="email" placeholder="email" v-model="email">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Signup</button>
    </form>
    <h2>{{ response }}</h2>
  </div>
</template>

<script>
import Error from '@/components/Error.vue';

export default {
    name: 'Signup',
    data() {
        return {
            username: '',
            email: '',
            password: '',
            response: '',
        };
    },
    components: {
        Error,
    },
    methods: {
        signup() {
            /* eslint-disable no-param-reassign */
            this.$el.querySelectorAll('.error').forEach((c) => {
                c.style.display = 'none';
            });

            const user = {
                username: this.username,
                email: this.email,
                password: this.password,
            };
            this.$http.post('/api/user', user)
                .then(() => {
                    this.$emit('signup');
                }).catch((e) => {
                    if (e.response.data.kind === this.UNIQUENESS_VIOLATION) {
                        this.$el.querySelector('.uniqueness.error').style.display = 'inline-block';
                    } else if (e.response.data.kind === this.DATA_VIOLATION) {
                        this.$el.querySelector('.format.error').style.display = 'inline-block';
                    } else {
                        console.log(`Error: ${e.response.data.message}`);
                        this.$el.querySelector('.unknown.error').style.display = 'inline-block';
                    }
                });
            /* eslint-enable no-param-reassign */
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.signup {
  position: relative;
}

input {
  margin-right: 0;
  margin-left: 0;
}

form {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.error {
  display: none;
  position: absolute;
  width: 180px;
  right: 180px;
}
</style>
