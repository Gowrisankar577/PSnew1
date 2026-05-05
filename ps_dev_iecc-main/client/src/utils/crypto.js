import CryptoJS from "crypto-js";

const secretKey = process.env.REACT_APP_KEY;

const encryptCourseId = (value) => {
  return CryptoJS.AES.encrypt(value.toString(), secretKey).toString();
};

const decryptCourseId = (value) => {
  try {
    const bytes = CryptoJS.AES.decrypt(value, secretKey);
    return bytes.toString(CryptoJS.enc.Utf8);
  } catch (error) {
    return "";
  }
};

export { encryptCourseId, decryptCourseId };
