<template>
  <div class="changepassword">
    <form @submit.prevent="changePassword">
      <input type="password" placeholder="password" v-model="password">
      <input type="password" placeholder="new password" v-model="newPassword">
      <button type="submit">Change Password</button>
    </form>
    <Error class="error empty">
      Please enter your password along with a new password
    </Error>
    <Error class="error login">
      Password incorrect. <a href="todo">Forgot password?</a>
    </Error>
    <Error class="error unknown">
      Could not change password, please try again later
    </Error>
  </div>
</template>

<script>
import Error from '@/components/Error.vue';

export default {
    name: 'Login',
    data() {
        return {
            password: '',
            newPassword: '',
        };
    },
    components: {
        Error,
    },
    methods: {
        changePassword() {
            /* eslint-disable no-param-reassign */
            this.$el.querySelectorAll('.error').forEach((c) => {
                c.style.display = 'none';
            });

            if (this.password === '' || this.newPassword === '') {
                this.$el.querySelector('.empty.error').style.display = 'inline-block';
                return;
            }

            const passwordChange = {
                password: this.password,
                newPassword: this.newPassword,
            };

            this.$http.post('/api/user/password', passwordChange)
                .then(() => {
                    this.password = '';
                    this.newPassword = '';
                    alert('Successfully changed your password'); // eslint-disable-line no-alert
                })
                .catch((e) => {
                    if (e.response.status === 401) {
                        this.$el.querySelector('.login.error').style.display = 'inline-block';
                    } else {
                        console.log(e.response);
                        this.$el.querySelector('.unknown.error').style.display = 'inline-block';
                        if (e) {
                            console.log(`${JSON.stringify(e)}`);
                        }
                    }
                });
            /* eslint-enable no-param-reassign */
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.changepassword {
  position: relative;
}

input {
  margin-top: 0;
  margin-bottom: 0;
  width: 120px;
}

button {
  margin: 0 $margin0;
}

form {
  display: flex;
}

.error {
  position: absolute;
  display: none;
}
</style>
