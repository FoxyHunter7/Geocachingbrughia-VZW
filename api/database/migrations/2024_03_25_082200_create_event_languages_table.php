GNU nano 7.2                        ../geocachingbrughia-vzw-old/database/migrations/2024_03_25_082200_create_event_languages_table.php                                 <?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('event_languages', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->foreignId('event_id')->constrained()->onUpdate('cascade')->onDelete('cascade');
            $table->string('lang_code');
            $table->foreign('lang_code')->references('code')->on('languages')->onUpdate('cascade')->onDelete('cascade');
            $table->unique(['event_id','lang_code']);

            $table->text('description');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('event_languages');
    }
};
