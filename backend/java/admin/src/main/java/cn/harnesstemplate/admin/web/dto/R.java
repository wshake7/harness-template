package cn.harnesstemplate.admin.web.dto;

public record R<T>(int code, String msg, T data) {

    public static final int CODE_SUCCESS = 1;
    public static final int CODE_FAIL = 2;
    public static final int CODE_REQUEST_EXPIRED = 3;
    public static final int CODE_LOGIN_FAIL = 100;
    public static final int CODE_AUTH_UNAUTHORIZED = 200;

    public static <T> R<T> ok(T data) {
        return new R<>(CODE_SUCCESS, "success", data);
    }

    public static <T> R<T> ok() {
        return new R<>(CODE_SUCCESS, "success", null);
    }

    public static <T> R<T> fail(String msg) {
        return new R<>(CODE_FAIL, msg, null);
    }

    public static <T> R<T> fail() {
        return new R<>(CODE_FAIL, "fail", null);
    }

    public static <T> R<T> unauthorized() {
        return new R<>(CODE_AUTH_UNAUTHORIZED, "未授权", null);
    }
}
