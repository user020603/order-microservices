package com.example.user.service;

import com.example.user.dto.UserDTO;
import com.example.user.event.UserEvent;
import com.example.user.kafka.KafkaProducer;
import com.example.user.model.User;
import com.example.user.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final RedisTemplate<String, UserDTO> redisTemplate;
    private final KafkaProducer kafkaProducer;

    @Cacheable(value = "users", key = "#id")
    public UserDTO getUser(Long id) {
        return userRepository.findById(id)
                .map(this::convertToDTO)
                .orElseThrow(() -> new RuntimeException("User not found"));
    }

    public UserDTO createUser(UserDTO userDTO) {
        if (userRepository.existsByEmail(userDTO.getEmail())) {
            throw new RuntimeException("Email already exists");
        }

        User user = convertToEntity(userDTO);
        user = userRepository.save(user);

        UserEvent event = new UserEvent();
        event.setEventType("CREATED");
        event.setUserId(user.getId());
        event.setUsername(user.getUsername());
        event.setEmail(user.getEmail());
        kafkaProducer.sendUserEvent(event);
        System.out.println("Produce send message from Kafka");

        return convertToDTO(user);
    }

    @CacheEvict(value = "users", key = "#id")
    public UserDTO updateUser(Long id, UserDTO userDTO) {
        User user = userRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("User not found"));

        user.setUsername(userDTO.getUsername());
        user.setEmail(userDTO.getEmail());
        user = userRepository.save(user);
        return convertToDTO(user);
    }

    @CacheEvict(value = "users", key = "#id")
    public void deleteUser(Long id) {
        userRepository.deleteById(id);
    }

    private UserDTO convertToDTO(User user) {
        UserDTO dto = new UserDTO();
        dto.setId(user.getId());
        dto.setUsername(user.getUsername());
        dto.setEmail(user.getEmail());
        return dto;
    }

    private User convertToEntity(UserDTO dto) {
        User user = new User();
        user.setUsername(dto.getUsername());
        user.setEmail(dto.getEmail());
        user.setPassword(dto.getPassword());
        return user;
    }
}

