use rsa::{pkcs1::{DecodeRsaPrivateKey, DecodeRsaPublicKey}, RsaPrivateKey, RsaPublicKey, Pkcs1v15Encrypt};
use rand::{rngs::OsRng, RngCore};
use aes_gcm::{aead::Aead, Aes256Gcm, KeyInit, Nonce};
use std::error::Error;

type AesGcm = Aes256Gcm;

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

pub fn decrypt_config(ciphertext: &[u8], key: &[u8]) -> Result<Vec<u8>, Box<dyn Error>> {
    let mut prepared_key = key.to_vec();
    
    while prepared_key.len() < 32 {
        prepared_key.push(b'0');
    }
    
    if prepared_key.len() > 32 {
        prepared_key.truncate(32);
    }
    
    let cipher = AesGcm::new_from_slice(&prepared_key)?;
    
    let nonce_size = 12;
    
    if ciphertext.len() < nonce_size {
        return Err("Ciphertext too short".into());
    }
    
    let (nonce_bytes, encrypted_data) = ciphertext.split_at(nonce_size);
    let nonce = Nonce::from_slice(nonce_bytes);
    
    let plaintext = cipher.decrypt(nonce, encrypted_data).map_err(|_| "Failed to decrypt config")?;
    
    Ok(plaintext)
}

pub fn encrypt_config(config_data: &[u8], key: &[u8]) -> Result<Vec<u8>, Box<dyn Error>> {
    let mut prepared_key = key.to_vec();
    
    while prepared_key.len() < 32 {
        prepared_key.push(b'0');
    }
    
    if prepared_key.len() > 32 {
        prepared_key.truncate(32);
    }
    
    let cipher = AesGcm::new_from_slice(&prepared_key)?;
    
    let mut nonce_bytes = [0u8; 12];
    rand::thread_rng().fill_bytes(&mut nonce_bytes);
    let nonce = Nonce::from_slice(&nonce_bytes);
    
    let ciphertext = cipher.encrypt(nonce, config_data).map_err(|_| "Failed to encrypt config")?;
    
    let mut result = nonce_bytes.to_vec();
    result.extend_from_slice(&ciphertext);
    
    Ok(result)
}