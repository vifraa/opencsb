<script>

    async function handleLogin(event) {


        const formData = new FormData();
        formData.append('username', event.target.username.value)
        formData.append('password', event.target.password.value)

        await fetch('http://0.0.0.0:8080/login', {
            method: 'POST',
            body: formData
        })
            .then((res) => res.json())
            .then((response) => console.log(response)); 
            
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


