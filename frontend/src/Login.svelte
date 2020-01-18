<script>

    async function handleLogin(event) {
        console.log(event)
        console.log(event.target.username.value)
        console.log(event.target.password.value)
    }

</script>


<div class="login-container">
    <h1>Login</h1>
    <form on:submit|preventDefault="{handleLogin}"
        class="d-flex flex-column justify-content-center align-items-center">
        <input type="text" id="username" required placeholder="Username">
        <input type="password" id="password" required placeholder="Password">
        <button type="submit" class="btn btn-primary">Login</button>
    </form>
</div>

<style>
    .login-container {
        background-color: #ffffff;
        width: fit-content;
        padding: 30px;
        padding: 15 15 30 30;
        margin: auto;

    }

    h1 {
        font-size: 26px;
        text-transform: uppercase;
    }

    button {
        width: 100%;
    }

    .login-container>* , form>*{
        margin-bottom: 8px;
    }

</style>


