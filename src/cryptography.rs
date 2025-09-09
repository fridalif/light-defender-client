use rsa::{pkcs1::{DecodeRsaPrivateKey, DecodeRsaPublicKey, EncodeRsaPrivateKey, EncodeRsaPublicKey}, Pkcs1v15Encrypt, RsaPrivateKey, RsaPublicKey};
use rand::{rngs::OsRng, RngCore};
use aes_gcm::{aead::Aead, Aes256Gcm, KeyInit, Nonce};
use std::error::Error;
use base64::{Engine as _, engine::general_purpose::STANDARD};

type AesGcm = Aes256Gcm;

pub fn generate_rsa_keys() -> Result<(String, String), Box<dyn Error>> {
    let mut rng = OsRng;
    
    let private_key = RsaPrivateKey::new(&mut rng, 2048)?;
    
    let public_key = RsaPublicKey::from(&private_key);
    
    let private_key_der = private_key.to_pkcs1_der()?;
    let private_key_bytes = private_key_der.as_bytes();
    
    let public_key_der = public_key.to_pkcs1_der()?;
    let public_key_bytes = public_key_der.as_bytes();
    
    let private_key_base64 = STANDARD.encode(private_key_bytes);
    let public_key_base64 = STANDARD.encode(public_key_bytes);
    
    Ok((private_key_base64, public_key_base64))
}

pub fn base64_to_bytes(base64_string: &str) -> Result<Vec<u8>, Box<dyn Error>> {
    Ok(STANDARD.decode(base64_string)?)
}

pub fn bytes_to_base64(bytes: &[u8]) -> String {
    STANDARD.encode(bytes)
}



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