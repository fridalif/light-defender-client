use rsa::{pkcs1::{DecodeRsaPrivateKey, DecodeRsaPublicKey}, RsaPrivateKey, RsaPublicKey, Pkcs1v15Encrypt};
use rand::rngs::OsRng;

pub fn encrypt_for_client(data: &Vec<u8>, client_pk: &Vec<u8>) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let public_key = RsaPublicKey::from_pkcs1_der(client_pk)?;
    let mut rng = OsRng;
    let encrypted_bytes = public_key.encrypt(&mut rng, Pkcs1v15Encrypt, data)?;
    Ok(encrypted_bytes)
}

pub fn decrypt_from_client(data: &Vec<u8>, server_private_key: &Vec<u8>) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let private_key = RsaPrivateKey::from_pkcs1_der(server_private_key)?;
    let decrypted_bytes = private_key.decrypt(Pkcs1v15Encrypt, data)?;
    Ok(decrypted_bytes)
}