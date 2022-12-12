<script>
export default {
    data: function () {
        return {
            identifier: "",
            errormsg: null,
            data: null,
        }
    },
    methods: {
        async sendData() {
            try {
                let response = await this.$axios.post("/session", {
                    identifier: this.identifier,
                });
                this.handleResponse(response);
            } catch (error) {
                this.handleError(error);
            }
        },

        // Handle a good response from the server
        handleResponse(response) {

            // Check if the response is 200 (OK)
            if (response.status === 200) {
                // If the response is 200, then the user is logged in
                // We can redirect him to the home page
                sessionStorage.setItem("token", response.data.token);
                alert("Login successful, your token is: " + sessionStorage.getItem("token"));
                this.$router.push("/");
            }
        },

        // Handle an error from the server
        handleError(error) {

            // Create an alert with the error from the server
            alert(error.response.data.message);

            // Print on console the error from the server (for debugging)
            console.log(error);
            console.log(error.data)
            console.log(error.status)

        }

    },

}
</script>




<template>
    <div class="container mt-5">
        <div class="row">
            <div class="col-3"></div>
            <div class="col-6">
                <form>
                    <div class="form-group">
                        <label for="exampleDropdownFormEmail2">Username</label>
                        <input type="text" v-model="identifier" class="form-control" placeholder="username">
                    </div>
                    <button @click="sendData()" type="submit" class="btn btn-primary">Sign in</button>
                </form>
            </div>
            <div class="col-3"></div>
        </div>

    </div>


</template>