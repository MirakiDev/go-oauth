export const logout = async () => {
    await fetch("http://localhost:8080/auth/logout/google", {
        credentials: "include",
    });

    window.location.href = "http://localhost:5173"
}