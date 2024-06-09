<script lang="ts">
    import { json } from "@sveltejs/kit";


let title: String = ""
const handleSubmit = () => {
    console.log(title)
    fetch("/api/make-todo", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "title": title,
            "userId": 1
        })
    }).catch(err => console.error(err))
    title = ""
}

const handleTest = () => {
    fetch("/api/test-ping")
        .then(v => console.log(v))
        .catch(err => console.error(err))
}


</script>

<form>
    <label for="title">title: </label>
    <input type="text" bind:value={title}>
    <button on:click={handleSubmit}>OK!</button>
    <button on:click={handleTest}>test</button>
</form>


