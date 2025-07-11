from RealtimeSTT import AudioToTextRecorder
import time

def main():
    """Função principal para executar a transcrição com Whisper."""

    # --- Configurações ---
    # Modelo do Whisper a ser usado (tiny, base, small, medium, large).
    # Modelos maiores são mais precisos, mas exigem mais recursos.
    MODEL_SIZE = "medium"
    # Idioma da transcrição. "pt" para português.
    LANGUAGE = "pt"
    # Dispositivo de computação ("cpu" ou "cuda" para GPU NVIDIA).
    DEVICE = "cuda" 
    # --- Fim das Configurações ---

    # Função de callback que será chamada com o texto transcrito
    def process_text(text):
        print("Transcrito:", text)

    print("Inicializando o gravador com o modelo Whisper...")
    
    try:
        # Cria a instância do gravador com as configurações desejadas
        recorder = AudioToTextRecorder(
            model=MODEL_SIZE,
            language=LANGUAGE,
            device=DEVICE,
            spinner=True # Desativa a animação de "loading" no console
        )

        print("Gravador inicializado. Fale no microfone (o silêncio encerra a gravação de um segmento).")
        print("Pressione Ctrl+C para parar.")

        # Loop infinito para processar o texto continuamente usando o callback
        while True:
            recorder.text(process_text)

    except KeyboardInterrupt:
        print("\nFinalizando a transcrição.")
    except Exception as e:
        print(f"Ocorreu um erro: {type(e).__name__}: {e}")

if __name__ == '__main__':
    main()